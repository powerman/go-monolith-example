package jsonrpc2

import (
	"github.com/prometheus/client_golang/prometheus"

	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/pkg/jsonrpc2x"
)

//nolint:gochecknoglobals // By design.
var (
	metric jsonrpc2x.Metrics
)

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(reg *prometheus.Registry) {
	metric = jsonrpc2x.NewMetrics(reg, app.ServiceName, "jsonrpc2",
		map[string]interface{}{
			"RPC": new(Server),
		},
		api.ErrsCommon,
		api.ErrsExtra,
	)
}
