package def

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics provides common metrics used by all packages.
type Metrics struct {
	PanicsTotal prometheus.Counter
}

// NewMetrics registers and returns common metrics used by all packages
// (subsystems) of given service (namespace).
func NewMetrics(reg *prometheus.Registry, service string) Metrics {
	var metrics Metrics
	metrics.PanicsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: service,
			Name:      "panics_total",
			Help:      "Amount of recovered panics.",
		},
	)
	reg.MustRegister(metrics.PanicsTotal)
	return metrics
}
