package auth

import (
	"testing"

	"github.com/powerman/check"
	_ "github.com/smartystreets/goconvey/convey"

	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/ms/auth/internal/config"
	"github.com/powerman/go-monolith-example/ms/auth/internal/dal"
	"github.com/powerman/go-monolith-example/ms/auth/internal/srv/grpc"
	"github.com/powerman/go-monolith-example/pkg/def"
)

func TestMain(m *testing.M) {
	def.Init()
	dal.InitMetrics(reg)
	app.InitMetrics(reg)
	grpc.InitMetrics(reg)
	cfg = config.MustGetServeTest()
	check.TestMain(m)
}

type tLogger check.C

func (l tLogger) Print(v ...interface{}) { l.Log(v...) }

// Const shared by tests. Recommended naming scheme: <dataType><Variant>.
var (
	cfg *config.ServeConfig
	ctx = def.NewContext(app.ServiceName)
)
