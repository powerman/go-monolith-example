package app

import (
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/prometheus/client_golang/prometheus"
)

//nolint:gochecknoglobals // By design.
var (
	Metric def.Metrics // Common metrics used by all packages.
	metric struct {
		ErrAccessDeniedTotal prometheus.Counter
	}
)

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(reg *prometheus.Registry) {
	const subsystem = "app"

	Metric = def.NewMetrics(reg, ServiceName)

	metric.ErrAccessDeniedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "ErrAccessDenied_total",
			Help:      "Amount of Access Denied errors.",
		},
	)
	reg.MustRegister(metric.ErrAccessDeniedTotal)
}
