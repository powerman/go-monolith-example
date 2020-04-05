package main

import (
	"regexp"

	"github.com/powerman/go-monolith-example/internal/concurrent"
	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/serve"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // Flags and metrics are service-global anyway.
var (
	reg = prometheus.DefaultRegisterer.(*prometheus.Registry)
	cfg struct {
		metrics config.Metrics
	}
)

// Service implements main.embeddedService interface.
type Service struct{}

// Name implements main.embeddedService interface.
func (Service) Name() string { return "main" }

// Init implements main.embeddedService interface.
func (Service) Init(_, serveCmd *cobra.Command) {
	namespace := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(exe, "_")
	InitMetrics(namespace)

	cfg.metrics.AddTo(serveCmd, "", def.MetricsPort)
}

// Serve implements main.embeddedService interface.
func (Service) Serve(_, ctxShutdown Ctx, shutdown func()) error {
	return concurrent.Serve(ctxShutdown, shutdown,
		serveMetrics)
}

func serveMetrics(ctx Ctx) error {
	return serve.Metrics(ctx, cfg.metrics, reg)
}
