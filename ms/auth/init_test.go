package auth

import (
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/ms/auth/internal/config"
	"github.com/powerman/go-monolith-example/ms/auth/internal/srv/grpc"
	"github.com/powerman/go-monolith-example/pkg/def"
)

func TestMain(m *testing.M) {
	def.Init()
	// dal.InitMetrics(reg) // TODO
	app.InitMetrics(reg)
	grpc.InitMetrics(reg)
	cfg = config.MustGetServeTest()
	check.TestMain(m)
}

// Const shared by tests. Recommended naming scheme: <dataType><Variant>.
var (
	cfg *config.ServeConfig
	ctx = def.NewContext(app.ServiceName)
)
