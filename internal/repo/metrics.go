package repo

import (
	"time"

	"github.com/powerman/go-monolith-example/internal/reflectx"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics contains general metrics for DAL methods.
type Metrics struct {
	callTotal    *prometheus.CounterVec
	callErrTotal *prometheus.CounterVec
	callDuration *prometheus.HistogramVec
}

const methodLabel = "method"

// NewMetrics registers and returns common DAL metrics used by all
// services (namespace).
func NewMetrics(reg *prometheus.Registry, service string, methodsFrom interface{}) (metric Metrics) {
	const subsystem = "dal"

	metric.callTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: service,
			Subsystem: subsystem,
			Name:      "call_total",
			Help:      "Amount of DAL calls.",
		},
		[]string{methodLabel},
	)
	reg.MustRegister(metric.callTotal)
	metric.callErrTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: service,
			Subsystem: subsystem,
			Name:      "errors_total",
			Help:      "Amount of DAL errors.",
		},
		[]string{methodLabel},
	)
	reg.MustRegister(metric.callErrTotal)
	metric.callDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: service,
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
		metric.callTotal.With(l)
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
			m.callTotal.With(l).Inc()
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
