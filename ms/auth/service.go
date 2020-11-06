// Package auth provides embedded microservice.
package auth

import (
	"context"
	"crypto/tls"

	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	grpcpkg "google.golang.org/grpc"

	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/ms/auth/internal/config"
	"github.com/powerman/go-monolith-example/ms/auth/internal/dal"
	"github.com/powerman/go-monolith-example/ms/auth/internal/srv/grpc"
	"github.com/powerman/go-monolith-example/pkg/concurrent"
	"github.com/powerman/go-monolith-example/pkg/serve"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

var reg = prometheus.NewPedanticRegistry() //nolint:gochecknoglobals // Metrics are global anyway.

// Service implements main.embeddedService interface.
type Service struct {
	cfg     *config.ServeConfig
	cert    tls.Certificate
	certInt tls.Certificate
	repo    *dal.Repo
	appl    *app.App
	srv     *grpcpkg.Server
	srvInt  *grpcpkg.Server
}

// Name implements main.embeddedService interface.
func (s *Service) Name() string { return app.ServiceName }

// Init implements main.embeddedService interface.
func (s *Service) Init(sharedCfg *config.SharedCfg, cmd, serveCmd *cobra.Command) error {
	// dal.InitMetrics(reg) TODO
	app.InitMetrics(reg)
	grpc.InitMetrics(reg)

	return config.Init(sharedCfg, config.FlagSets{
		Serve: serveCmd.Flags(),
	})
}

// RunServe implements main.embeddedService interface.
func (s *Service) RunServe(ctxStartup, ctxShutdown Ctx, shutdown func()) (err error) {
	log := structlog.FromContext(ctxShutdown, nil)
	if s.cfg == nil {
		s.cfg, err = config.GetServe()
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
		s.appl = app.New(s.repo, app.Config{})
	}

	s.srv = grpc.NewServer(s.appl, grpc.Config{
		Cert: &s.cert,
	})
	s.srvInt = grpc.NewServerInt(s.appl, grpc.Config{
		Cert: &s.certInt,
	})

	err = concurrent.Serve(ctxShutdown, shutdown,
		s.serveMetrics,
		s.serveGRPC,
		s.serveGRPCInt,
	)
	if err != nil {
		return log.Err("failed to serve", "err", err)
	}
	return nil
}

func (s *Service) connectRepo(ctx Ctx) (interface{}, error) {
	return dal.New(), nil
}

func (s *Service) serveMetrics(ctx Ctx) error {
	return serve.Metrics(ctx, s.cfg.MetricsAddr, reg)
}

func (s *Service) serveGRPC(ctx Ctx) error {
	return serve.GRPC(ctx, s.cfg.Addr, s.srv, "gRPC")
}

func (s *Service) serveGRPCInt(ctx Ctx) error {
	return serve.GRPC(ctx, s.cfg.AddrInt, s.srvInt, "gRPC internal")
}
