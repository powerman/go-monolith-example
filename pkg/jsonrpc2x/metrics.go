package jsonrpc2x

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/powerman/go-monolith-example/pkg/reflectx"
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
	failedLabel = "failed"
)

// NewMetrics registers and returns common JSON-RPC 2.0 metrics used by
// all services (namespace).
func NewMetrics( //nolint:funlen // By design.
	reg *prometheus.Registry,
	service string,
	subsystem string,
	methodsFrom map[string]interface{},
	errsCommon []error,
	errsExtra map[string][]error,
) (
	metric Metrics,
) {
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
		[]string{methodLabel, failedLabel},
	)
	reg.MustRegister(metric.reqDuration)

	commonCodes := []string{""} // Successful RPC.
	commonCodes = append(commonCodes, codes(errsCommon)...)
	for name, rcvr := range methodsFrom {
		for _, methodName := range reflectx.RPCMethodsOf(rcvr) {
			methodName = name + "." + methodName
			if _, ok := errsExtra[methodName]; !ok {
				panic(fmt.Sprintf("missing ErrsExtra[%s]", methodName))
			}
			codes := append(commonCodes, codes(errsExtra[methodName])...) //nolint:gocritic // Not same slice.
			for _, code := range codes {
				l := prometheus.Labels{
					methodLabel: methodName,
					codeLabel:   code,
				}
				metric.reqTotal.With(l)
			}
			for _, failed := range []string{"true", "false"} {
				l := prometheus.Labels{
					methodLabel: methodName,
					failedLabel: failed,
				}
				metric.reqDuration.With(l)
			}
		}
	}

	return metric
}
