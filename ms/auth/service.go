// Package auth provides embedded microservice.
package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	grpcpkg "google.golang.org/grpc"

	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/ms/auth/internal/config"
	"github.com/powerman/go-monolith-example/ms/auth/internal/dal"
	"github.com/powerman/go-monolith-example/ms/auth/internal/migrations"
	"github.com/powerman/go-monolith-example/ms/auth/internal/srv/grpc"
	"github.com/powerman/go-monolith-example/ms/auth/internal/srv/grpcgw"
	"github.com/powerman/go-monolith-example/pkg/cobrax"
	"github.com/powerman/go-monolith-example/pkg/concurrent"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/netx"
	"github.com/powerman/go-monolith-example/pkg/serve"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

var reg = prometheus.NewPedanticRegistry() //nolint:gochecknoglobals // Metrics are global anyway.

// Service implements main.embeddedService interface.
type Service struct {
	cfg     *config.ServeConfig
	ca      *x509.CertPool
	cert    tls.Certificate
	certInt tls.Certificate
	repo    *dal.Repo
	appl    *app.App
	srv     *grpcpkg.Server
	srvInt  *grpcpkg.Server
	mux     *http.ServeMux
}

// Name implements main.embeddedService interface.
func (s *Service) Name() string { return app.ServiceName }

// Init implements main.embeddedService interface.
func (s *Service) Init(sharedCfg *config.SharedCfg, cmd, serveCmd *cobra.Command) error {
	dal.InitMetrics(reg)
	app.InitMetrics(reg)
	grpc.InitMetrics(reg)

	ctx := def.NewContext(app.ServiceName)
	goosePostgresCmd := cobrax.NewGoosePostgresCmd(ctx, migrations.Goose(), config.GetGoosePostgres)
	cmd.AddCommand(goosePostgresCmd)

	return config.Init(sharedCfg, config.FlagSets{
		Serve:         serveCmd.Flags(),
		GoosePostgres: goosePostgresCmd.Flags(),
	})
}

// RunServe implements main.embeddedService interface.
func (s *Service) RunServe(ctxStartup, ctxShutdown Ctx, shutdown func()) (err error) {
	log := structlog.FromContext(ctxShutdown, nil)
	if s.cfg == nil {
		s.cfg, err = config.GetServe()
	}
	if err == nil {
		s.ca, err = netx.LoadCACert(s.cfg.TLSCACert)
	}
	if err == nil {
		s.cert, err = tls.LoadX509KeyPair(s.cfg.TLSCert, s.cfg.TLSKey)
	}
	if err == nil {
		s.certInt, err = tls.LoadX509KeyPair(s.cfg.TLSCertInt, s.cfg.TLSKeyInt)
	}
	if err != nil {
		return log.Err("failed to get config", "err", err)
	}

	err = concurrent.Setup(ctxStartup, map[interface{}]concurrent.SetupFunc{
		&s.repo: s.connectRepo,
	})
	if err != nil {
		return log.Err("failed to connect", "err", err)
	}

	if s.appl == nil {
		s.appl = app.New(s.repo, app.Config{
			Secret: s.cfg.Secret,
		})
	}

	s.srv = grpc.NewServer(s.appl, grpc.Config{
		CtxShutdown: ctxShutdown,
		Cert:        &s.cert,
	})
	s.srvInt = grpc.NewServerInt(s.appl, grpc.Config{
		CtxShutdown: ctxShutdown,
		Cert:        &s.certInt,
	})
	s.mux, err = grpcgw.NewServer(grpcgw.Config{
		CtxShutdown:      ctxShutdown,
		Endpoint:         s.cfg.AuthAddr,
		CA:               s.ca,
		GRPCGWPattern:    "/",
		OpenAPIPattern:   "/openapi/", // Also hardcoded in web/static/swagger-ui/index.html.
		SwaggerUIPattern: "/swagger-ui/",
	})
	if err != nil {
		return log.Err("failed to setup grpc-gateway", "err", err)
	}

	err = concurrent.Serve(ctxShutdown, shutdown,
		s.serveMetrics,
		s.serveGRPC,
		s.serveGRPCInt,
		s.serveGRPCGW,
	)
	if err != nil {
		return log.Err("failed to serve", "err", err)
	}
	return nil
}

func (s *Service) connectRepo(ctx Ctx) (interface{}, error) {
	return dal.New(ctx, s.cfg.GoosePostgresDir, s.cfg.Postgres)
}

func (s *Service) serveMetrics(ctx Ctx) error {
	return serve.Metrics(ctx, s.cfg.BindMetricsAddr, reg)
}

func (s *Service) serveGRPC(ctx Ctx) error {
	return serve.GRPC(ctx, s.cfg.BindAddr, s.srv, "gRPC")
}

func (s *Service) serveGRPCInt(ctx Ctx) error {
	return serve.GRPC(ctx, s.cfg.BindAddrInt, s.srvInt, "gRPC internal")
}

func (s *Service) serveGRPCGW(ctx Ctx) error {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{s.cert},
		MinVersion:   tls.VersionTLS12,
	}
	return serve.HTTP(ctx, s.cfg.BindGRPCGWAddr, tlsConfig, s.mux, "grpc-gateway")
}
