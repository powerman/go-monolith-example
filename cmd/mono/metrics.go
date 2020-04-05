package main

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(namespace string) {
	version := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "build_info",
			Help:      "A metric with a constant '1' value labeled by build-time details.",
		},
		[]string{"version", "goversion"},
	)
	prometheus.MustRegister(version)

	version.With(prometheus.Labels{
		"version":   ver,
		"goversion": runtime.Version(),
	}).Set(1)
}
