// Package example provides embedded microservice.
package example

import (
	"context"

	"github.com/powerman/go-monolith-example/internal/apiauth"
	"github.com/powerman/go-monolith-example/internal/cobrax"
	"github.com/powerman/go-monolith-example/internal/concurrent"
	"github.com/powerman/go-monolith-example/internal/event"
	"github.com/powerman/go-monolith-example/internal/serve"
	"github.com/powerman/go-monolith-example/ms/example/internal/api"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/ms/example/internal/config"
	"github.com/powerman/go-monolith-example/ms/example/internal/dal"
	"github.com/powerman/go-monolith-example/ms/example/internal/migrations"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

//nolint:gochecknoglobals // Flags and metrics are service-global anyway.
var (
	reg = prometheus.NewPedanticRegistry()
	cfg *config.ServeConfig
)

// Service implements main.embeddedService interface.
type Service struct{}

// Name implements main.embeddedService interface.
func (Service) Name() string { return app.ServiceName }

// Init implements main.embeddedService interface.
func (Service) Init(genericCfg *config.GenericCfg, cmd, serveCmd *cobra.Command) error {
	dal.InitMetrics(reg)
	app.InitMetrics(reg)
	api.InitMetrics(reg)

	gooseCmd := cobrax.NewGooseCmd(app.ServiceName, migrations.Goose(), config.GetGoose)
	cmd.AddCommand(gooseCmd)

	return config.Init(genericCfg, config.FlagSets{
		Serve: serveCmd.Flags(),
		Goose: gooseCmd.Flags(),
	})
}

//nolint:gochecknoglobals // For tests.
var (
	natsConn *event.NATSConn
	stanConn *event.STANConn
	repo     *dal.Repo
	authn    apiauth.Authenticator
	a        app.Appl
)

// Serve implements main.embeddedService interface.
func (Service) Serve(ctxStartup, ctxShutdown Ctx, shutdown func()) (err error) {
	log := structlog.FromContext(ctxShutdown, nil)
	if cfg == nil {
		cfg, err = config.GetServe()
	}
	if err != nil {
		return log.Err("failed to get config", "err", err)
	}

	err = concurrent.Setup(ctxStartup, map[interface{}]concurrent.SetupFunc{
		&natsConn: connectNATS,
		&repo:     connectRepo,
		&authn:    setupAuthn,
	})
	if err == nil && stanConn == nil {
		stanConn, err = event.ConnectSTAN(ctxStartup, cfg.STANClusterID, app.ServiceName, natsConn)
	}
	if natsConn != nil {
		defer log.WarnIfFail(natsConn.Drain)
	}
	if stanConn != nil {
		defer log.WarnIfFail(stanConn.Close)
	}
	if err != nil {
		return log.Err("failed to connect", "err", err)
	}

	if a == nil {
		a = app.New(repo, app.Config{})
	}

	err = concurrent.Serve(ctxShutdown, shutdown,
		natsConn.Monitor,
		stanConn.Monitor,
		serveMetrics,
		serveRPC,
	)
	if err != nil {
		return log.Err("failed to serve", "err", err)
	}
	return nil
}

func connectNATS(ctx Ctx) (interface{}, error) {
	return event.ConnectNATS(ctx, cfg.NATSUrls, app.ServiceName)
}

func connectRepo(ctx Ctx) (interface{}, error) {
	return dal.New(ctx, cfg.GooseDir, cfg.MySQLConfig)
}

func setupAuthn(_ Ctx) (interface{}, error) {
	return apiauth.NewAccessTokenParser(), nil
}

func serveMetrics(ctx Ctx) error {
	return serve.Metrics(ctx, cfg.MetricsAddr, reg)
}

func serveRPC(ctx Ctx) error {
	return serve.RPC(ctx, cfg.RPCAddr, api.New(a, authn))
}
