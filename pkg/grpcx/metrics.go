package grpcx

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

// NewServerMetrics returns gRPC server metrics.
// Do not forget to call .InitializeMetrics(server) on returned value.
func NewServerMetrics(reg *prometheus.Registry, namespace, subsystem string) *grpc_prometheus.ServerMetrics {
	serverMetrics := grpc_prometheus.NewServerMetrics(func(o *prometheus.CounterOpts) {
		o.Namespace = namespace
		o.Subsystem = subsystem
	})
	serverMetrics.EnableHandlingTimeHistogram(func(o *prometheus.HistogramOpts) {
		o.Namespace = namespace
		o.Subsystem = subsystem
	})
	reg.MustRegister(serverMetrics)
	return serverMetrics
}

// NewClientMetrics returns gRPC client metrics.
func NewClientMetrics(reg *prometheus.Registry, namespace, subsystem string) *grpc_prometheus.ClientMetrics {
	clientMetrics := grpc_prometheus.NewClientMetrics(func(o *prometheus.CounterOpts) {
		o.Namespace = namespace
		o.Subsystem = subsystem
	})
	clientMetrics.EnableClientHandlingTimeHistogram(func(o *prometheus.HistogramOpts) {
		o.Namespace = namespace
		o.Subsystem = subsystem
	})
	reg.MustRegister(clientMetrics)
	return clientMetrics
}
