// Package metrics provides embedded microservice.
package metrics

import (
	"context"
	"regexp"
	"strconv"

	"github.com/powerman/appcfg"
	"github.com/powerman/go-monolith-example/internal/concurrent"
	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/netx"
	"github.com/powerman/go-monolith-example/internal/serve"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

//nolint:gochecknoglobals // Config and metrics are service-global anyway.
var (
	fs         *pflag.FlagSet
	genericCfg *config.Cfg
	svcCfg     = &struct {
		MetricsPort appcfg.Port `env:"METRICS_PORT"`
	}{
		MetricsPort: appcfg.MustPort(strconv.Itoa(def.MetricsPort)),
	}

	reg = prometheus.NewPedanticRegistry()
	cfg struct {
		metricsAddr netx.Addr
	}
)

// Service implements main.embeddedService interface.
type Service struct{}

// Name implements main.embeddedService interface.
func (Service) Name() string { return "metrics" }

// Init implements main.embeddedService interface.
func (Service) Init(generic *config.Cfg, _, serveCmd *cobra.Command) error {
	namespace := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(def.ProgName, "_")
	initMetrics(reg, namespace)

	fs, genericCfg = serveCmd.Flags(), generic
	fromEnv := appcfg.NewFromEnv(config.EnvPrefix)
	err := appcfg.ProvideStruct(svcCfg, fromEnv)
	config.AddFlag(fs, &genericCfg.MetricsHost, "metrics.host", "host to serve Prometheus metrics")
	config.AddFlag(fs, &svcCfg.MetricsPort, "metrics.port", "port to serve Prometheus metrics")
	return err
}

// Serve implements main.embeddedService interface.
func (Service) Serve(_, ctxShutdown Ctx, shutdown func()) (err error) {
	log := structlog.FromContext(ctxShutdown, nil)
	cfg.metricsAddr = netx.NewAddr(genericCfg.MetricsHost.Value(&err), svcCfg.MetricsPort.Value(&err))
	if err != nil {
		return log.Err("failed to get config", "err", appcfg.WrapPErr(err, fs, genericCfg, svcCfg))
	}

	err = concurrent.Serve(ctxShutdown, shutdown,
		serveMetrics)
	if err != nil {
		return log.Err("failed to serve", "err", err)
	}
	return nil
}

func serveMetrics(ctx Ctx) error {
	return serve.Metrics(ctx, cfg.metricsAddr, reg)
}
