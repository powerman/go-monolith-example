package app

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/powerman/go-monolith-example/pkg/def"
)

//nolint:gochecknoglobals // Metrics are global anyway.
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

	Metric = def.NewMetrics(reg)

	metric.ErrAccessDeniedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "err_access_denied_total",
			Help:      "Amount of Access Denied errors.",
		},
	)
	reg.MustRegister(metric.ErrAccessDeniedTotal)
}
