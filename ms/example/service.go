// Package example provides embedded microservice.
package example

import (
	"context"
	"fmt"

	"github.com/powerman/go-monolith-example/internal/concurrent"
	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/event"
	"github.com/powerman/go-monolith-example/internal/migrate"
	"github.com/powerman/go-monolith-example/internal/serve"
	"github.com/powerman/go-monolith-example/ms/example/internal/api"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
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
	cfg struct {
		natsUrls      config.NATSUrls      // Not used yet, just an example.
		stanClusterID config.STANClusterID // Not used yet, just an example.
		metrics       config.Metrics
		mysql         config.MySQL
		rpc           config.RPC
	}
)

// Service implements main.embeddedService interface.
type Service struct{}

// Name implements main.embeddedService interface.
func (Service) Name() string { return app.ServiceName }

// Init implements main.embeddedService interface.
func (Service) Init(cmd, serveCmd *cobra.Command) {
	dal.InitMetrics(reg)
	app.InitMetrics(reg)
	api.InitMetrics(reg)

	mysqlDef := config.MySQLDef{
		User:     def.ExampleDBUser,
		Pass:     def.ExampleDBPass,
		DBName:   def.ExampleDBName,
		GooseDir: def.ExampleGooseDir,
	}

	gooseCmd := config.NewGooseCmd(runGoose)
	cfg.mysql.AddTo(gooseCmd, app.ServiceName, mysqlDef)
	cmd.AddCommand(gooseCmd)

	cfg.metrics.AddTo(serveCmd, app.ServiceName, def.ExampleMetricsPort)
	cfg.mysql.AddTo(serveCmd, app.ServiceName, mysqlDef)
	cfg.rpc.AddTo(serveCmd, app.ServiceName, def.ExampleRPCPort)
	cfg.natsUrls.AddTo(serveCmd)
	cfg.stanClusterID.AddTo(serveCmd)
}

func runGoose(cmd string) error {
	ctx := def.NewContext(app.ServiceName)
	return migrate.Run(ctx, migrations.Goose(), cfg.mysql.GooseDir(), cmd, cfg.mysql.Config())
}

//nolint:gochecknoglobals // For tests.
var (
	natsConn *event.NATSConn
	stanConn *event.STANConn
	repo     *dal.Repo
	a        app.Appl
)

// Serve implements main.embeddedService interface.
func (Service) Serve(ctxSetup, ctxShutdown Ctx, shutdown func()) error {
	err := concurrent.Setup(ctxSetup, map[interface{}]concurrent.SetupFunc{
		&natsConn: connectNATS,
		&repo:     connectRepo,
	})
	if err == nil && stanConn == nil {
		stanConn, err = event.ConnectSTAN(ctxSetup, cfg.stanClusterID.String(),
			app.ServiceName, natsConn)
	}
	if err != nil {
		return fmt.Errorf("connect: %w", err)
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

	log := structlog.FromContext(ctxShutdown, nil)
	log.WarnIfFail(stanConn.Close)
	log.WarnIfFail(natsConn.Drain)
	return err
}

func connectNATS(ctx Ctx) (interface{}, error) {
	return event.ConnectNATS(ctx, cfg.natsUrls.String(), app.ServiceName)
}

func connectRepo(ctx Ctx) (interface{}, error) {
	return dal.New(ctx, cfg.mysql.GooseDir(), cfg.mysql.Config())
}

func serveMetrics(ctx Ctx) error {
	return serve.Metrics(ctx, cfg.metrics, reg)
}

func serveRPC(ctx Ctx) error {
	return serve.RPC(ctx, cfg.rpc, api.New(a))
}
