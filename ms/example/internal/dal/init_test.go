package dal

import (
	"testing"

	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/gotest/testinit"
	"github.com/prometheus/client_golang/prometheus"
)

func TestMain(m *testing.M) { testinit.Main(m) }

const (
	serialMain = iota
	serialIntegration
)

func init() { testinit.Setup(serialMain, setupMain) }

func setupMain() {
	reg := prometheus.NewPedanticRegistry()
	def.Init()
	app.InitMetrics(reg)
	InitMetrics(reg)
}

var ctx = def.NewContext(app.ServiceName)
