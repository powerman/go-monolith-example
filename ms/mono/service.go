// Package mono provides embedded microservice.
package mono

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
		Port appcfg.Port `env:"MONO_ADDR_PORT"`
	}{
		Port: appcfg.MustPort(strconv.Itoa(config.MonoPort)),
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
func (s *Service) Name() string { return "mono" }

// Init implements main.embeddedService interface.
func (s *Service) Init(sharedCfg *config.Shared, _, serveCmd *cobra.Command) error {
	namespace := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(def.ProgName, "_")
	initMetrics(reg, namespace)

	fs, shared = serveCmd.Flags(), sharedCfg
	fromEnv := appcfg.NewFromEnv(config.EnvPrefix)
	err := appcfg.ProvideStruct(own, fromEnv)
	pfx := s.Name() + "."
	appcfg.AddPFlag(fs, &shared.AddrHostInt, "host-int", "internal host to serve")
	appcfg.AddPFlag(fs, &own.Port, pfx+"port", "port to serve monolith introspection")
	return err
}

// RunServe implements main.embeddedService interface.
func (s *Service) RunServe(_, ctxShutdown Ctx, shutdown func()) (err error) {
	log := structlog.FromContext(ctxShutdown, nil)
	s.cfg.metricsAddr = netx.NewAddr(shared.AddrHostInt.Value(&err), own.Port.Value(&err))
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
