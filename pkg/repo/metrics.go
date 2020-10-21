package repo

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/powerman/go-monolith-example/pkg/reflectx"
)

// Metrics contains general metrics for DAL methods.
type Metrics struct {
	callErrTotal *prometheus.CounterVec
	callDuration *prometheus.HistogramVec
}

const methodLabel = "method"

// NewMetrics registers and returns common DAL metrics used by all
// services (namespace).
func NewMetrics(reg *prometheus.Registry, namespace, subsystem string, methodsFrom interface{}) (metric Metrics) {
	metric.callErrTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "errors_total",
			Help:      "Amount of DAL errors.",
		},
		[]string{methodLabel},
	)
	reg.MustRegister(metric.callErrTotal)
	metric.callDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "call_duration_seconds",
			Help:      "DAL call latency.",
		},
		[]string{methodLabel},
	)
	reg.MustRegister(metric.callDuration)

	for _, methodName := range reflectx.MethodsOf(methodsFrom) {
		l := prometheus.Labels{
			methodLabel: methodName,
		}
		metric.callErrTotal.With(l)
		metric.callDuration.With(l)
	}

	return metric
}

func (m Metrics) instrument(method string, f func() error) func() error {
	return func() (err error) {
		start := time.Now()
		l := prometheus.Labels{methodLabel: method}
		defer func() {
			m.callDuration.With(l).Observe(time.Since(start).Seconds())
			if err != nil {
				m.callErrTotal.With(l).Inc()
			} else if err := recover(); err != nil {
				m.callErrTotal.With(l).Inc()
				panic(err)
			}
		}()
		return f()
	}
}
