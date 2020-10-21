// Package metrics provides embedded microservice.
package metrics

import (
	"context"
	"regexp"
	"strconv"

	"github.com/powerman/appcfg"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/pkg/concurrent"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/netx"
	"github.com/powerman/go-monolith-example/pkg/serve"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

//nolint:gochecknoglobals // Config, flags and metrics are global anyway.
var (
	fs     *pflag.FlagSet
	shared *config.Shared
	own    = &struct {
		MetricsPort appcfg.Port `env:"METRICS_ADDR_PORT"`
	}{
		MetricsPort: appcfg.MustPort(strconv.Itoa(config.MetricsPort)),
	}

	reg = prometheus.NewPedanticRegistry()
)

// Service implements main.embeddedService interface.
type Service struct {
	cfg struct {
		metricsAddr netx.Addr
	}
}

// Name implements main.embeddedService interface.
func (s *Service) Name() string { return "metrics" }

// Init implements main.embeddedService interface.
func (s *Service) Init(sharedCfg *config.Shared, _, serveCmd *cobra.Command) error {
	namespace := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(def.ProgName, "_")
	initMetrics(reg, namespace)

	fs, shared = serveCmd.Flags(), sharedCfg
	fromEnv := appcfg.NewFromEnv(config.EnvPrefix)
	err := appcfg.ProvideStruct(own, fromEnv)
	appcfg.AddPFlag(fs, &shared.MonoMetricsAddrHost, "metrics.host", "host to serve Prometheus metrics")
	appcfg.AddPFlag(fs, &own.MetricsPort, "metrics.port", "port to serve Prometheus metrics")
	return err
}

// RunServe implements main.embeddedService interface.
func (s *Service) RunServe(_, ctxShutdown Ctx, shutdown func()) (err error) {
	log := structlog.FromContext(ctxShutdown, nil)
	s.cfg.metricsAddr = netx.NewAddr(shared.MonoMetricsAddrHost.Value(&err), own.MetricsPort.Value(&err))
	if err != nil {
		return log.Err("failed to get config", "err", appcfg.WrapPErr(err, fs, shared, own))
	}

	err = concurrent.Serve(ctxShutdown, shutdown,
		s.serveMetrics,
	)
	if err != nil {
		return log.Err("failed to serve", "err", err)
	}
	return nil
}

func (s *Service) serveMetrics(ctx Ctx) error {
	return serve.Metrics(ctx, s.cfg.metricsAddr, reg)
}
