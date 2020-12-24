// Package example provides embedded microservice.
package example

import (
	"context"
	"crypto/x509"
	"net/http"

	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"

	"github.com/powerman/go-monolith-example/internal/apix"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/ms/example/internal/config"
	"github.com/powerman/go-monolith-example/ms/example/internal/dal"
	"github.com/powerman/go-monolith-example/ms/example/internal/migrations"
	"github.com/powerman/go-monolith-example/ms/example/internal/srv/jsonrpc2"
	"github.com/powerman/go-monolith-example/pkg/cobrax"
	"github.com/powerman/go-monolith-example/pkg/concurrent"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/natsx"
	"github.com/powerman/go-monolith-example/pkg/netx"
	"github.com/powerman/go-monolith-example/pkg/serve"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

var reg = prometheus.NewPedanticRegistry() //nolint:gochecknoglobals // Metrics are global anyway.

// Service implements main.embeddedService interface.
type Service struct {
	cfg      *config.ServeConfig
	ca       *x509.CertPool
	natsConn *natsx.NATSConn
	stanConn *natsx.STANConn
	repo     *dal.Repo
	authn    apix.Authn
	appl     *app.App
	mux      *http.ServeMux
}

// Name implements main.embeddedService interface.
func (s *Service) Name() string { return app.ServiceName }

// Init implements main.embeddedService interface.
func (s *Service) Init(sharedCfg *config.SharedCfg, cmd, serveCmd *cobra.Command) error {
	dal.InitMetrics(reg)
	app.InitMetrics(reg)
	jsonrpc2.InitMetrics(reg)

	ctx := def.NewContext(app.ServiceName)
	gooseMySQLCmd := cobrax.NewGooseMySQLCmd(ctx, migrations.Goose(), config.GetGooseMySQL)
	cmd.AddCommand(gooseMySQLCmd)

	return config.Init(sharedCfg, config.FlagSets{
		Serve:      serveCmd.Flags(),
		GooseMySQL: gooseMySQLCmd.Flags(),
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
	if err != nil {
		return log.Err("failed to get config", "err", err)
	}

	err = concurrent.Setup(ctxStartup, map[interface{}]concurrent.SetupFunc{
		&s.natsConn: s.connectNATS,
		&s.repo:     s.connectRepo,
		&s.authn:    s.setupAuthn,
	})
	if err == nil && s.stanConn == nil {
		s.stanConn, err = natsx.ConnectSTAN(ctxStartup, s.cfg.STANClusterID, app.ServiceName, s.natsConn)
	}
	if s.natsConn != nil {
		defer log.WarnIfFail(s.natsConn.Drain)
	}
	if s.stanConn != nil {
		defer log.WarnIfFail(s.stanConn.Close)
	}
	if err != nil {
		return log.Err("failed to connect", "err", err)
	}

	if s.appl == nil {
		s.appl = app.New(s.repo, app.Config{})
	}

	s.mux = jsonrpc2.NewServer(s.appl, s.authn, jsonrpc2.Config{
		Pattern: s.cfg.Path,
	})

	err = concurrent.Serve(ctxShutdown, shutdown,
		s.natsConn.Monitor,
		s.stanConn.Monitor,
		s.serveMetrics,
		s.serveHTTP,
	)
	if err != nil {
		return log.Err("failed to serve", "err", err)
	}
	return nil
}

func (s *Service) connectNATS(ctx Ctx) (interface{}, error) {
	return natsx.ConnectNATS(ctx, s.cfg.NATSURLs, app.ServiceName)
}

func (s *Service) connectRepo(ctx Ctx) (interface{}, error) {
	return dal.New(ctx, s.cfg.GooseMySQLDir, s.cfg.MySQL)
}

func (s *Service) setupAuthn(ctx Ctx) (interface{}, error) {
	return apix.NewAuthnClient(ctx, reg, app.ServiceName, s.ca, s.cfg.AuthAddrInt.String())
}

func (s *Service) serveMetrics(ctx Ctx) error {
	return serve.Metrics(ctx, s.cfg.BindMetricsAddr, reg)
}

func (s *Service) serveHTTP(ctx Ctx) error {
	return serve.HTTP(ctx, s.cfg.BindAddr, nil, s.mux, "JSON-RPC 2.0")
}
