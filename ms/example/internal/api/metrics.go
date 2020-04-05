package api

import (
	"github.com/powerman/go-monolith-example/internal/jsonrpc2x"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/prometheus/client_golang/prometheus"
)

//nolint:gochecknoglobals // By design.
var metric jsonrpc2x.Metrics

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(reg *prometheus.Registry) {
	metric = jsonrpc2x.NewMetrics(reg, app.ServiceName, new(API))
}
