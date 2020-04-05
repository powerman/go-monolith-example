package jsonrpc2x

import (
	"fmt"

	"github.com/powerman/go-monolith-example/internal/reflectx"
	"github.com/powerman/go-monolith-example/proto/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics contains general metrics for JSON-RPC 2.0 methods.
type Metrics struct {
	reqInFlight prometheus.Gauge
	reqTotal    *prometheus.CounterVec
	reqDuration *prometheus.HistogramVec
}

const (
	methodLabel = "method"
	codeLabel   = "code"
)

// NewMetrics registers and returns common JSON-RPC 2.0 metrics used by
// all services (namespace).
func NewMetrics(reg *prometheus.Registry, service string, methodsFrom interface{}) (metric Metrics) {
	const subsystem = "api"

	metric.reqInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: service,
			Subsystem: subsystem,
			Name:      "http_requests_in_flight",
			Help:      "Amount of currently processing API requests.",
		},
	)
	reg.MustRegister(metric.reqInFlight)
	metric.reqTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: service,
			Subsystem: subsystem,
			Name:      "http_requests_total",
			Help:      "Amount of processed API requests.",
		},
		[]string{methodLabel, codeLabel},
	)
	reg.MustRegister(metric.reqTotal)
	metric.reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: service,
			Subsystem: subsystem,
			Name:      "http_request_duration_seconds",
			Help:      "API request latency distributions.",
		},
		[]string{methodLabel, codeLabel},
	)
	reg.MustRegister(metric.reqDuration)

	commonCodes := []string{""} // Successful RPC.
	commonCodes = append(commonCodes, codes(rpc.ErrsCommon)...)
	for _, methodName := range reflectx.RPCMethodsOf(methodsFrom) {
		if _, ok := rpc.ErrsExtra[methodName]; !ok {
			panic(fmt.Sprintf("missing ErrsExtra[%s]", methodName))
		}
		codes := append(commonCodes, codes(rpc.ErrsExtra[methodName])...)
		for _, code := range codes {
			l := prometheus.Labels{
				methodLabel: methodName,
				codeLabel:   code,
			}
			metric.reqTotal.With(l)
			metric.reqDuration.With(l)
		}
	}

	return metric
}
