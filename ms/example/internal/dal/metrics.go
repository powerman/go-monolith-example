package dal

import (
	"github.com/powerman/go-monolith-example/internal/repo"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/prometheus/client_golang/prometheus"
)

//nolint:gochecknoglobals // By design.
var metric repo.Metrics

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(reg *prometheus.Registry) {
	metric = repo.NewMetrics(reg, app.ServiceName, new(app.Repo))
}
