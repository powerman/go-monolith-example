package grpc

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/pkg/grpcx"
)

// Metric contains general metrics for gRPC methods.
var metric struct { //nolint:gochecknoglobals // Metrics are global anyway.
	server *grpc_prometheus.ServerMetrics
}

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(reg *prometheus.Registry) {
	const subsystem = "grpc"
	metric.server = grpcx.NewServerMetrics(reg, app.ServiceName, subsystem)
}
