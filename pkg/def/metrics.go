package def

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics shared by all packages.
type Metrics struct {
	PanicsTotal           prometheus.Counter
	MisconfigurationTotal prometheus.Counter
}

// NewMetrics registers and returns metrics shared by all packages.
func NewMetrics(reg *prometheus.Registry) Metrics {
	var metrics Metrics
	metrics.PanicsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "panics_total",
			Help: "Amount of recovered panics.",
		},
	)
	reg.MustRegister(metrics.PanicsTotal)
	metrics.MisconfigurationTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "misconfiguration_total",
			Help: "Amount of failures because of incorrect configuration.",
		},
	)
	reg.MustRegister(metrics.MisconfigurationTotal)
	return metrics
}
