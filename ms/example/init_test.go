package example

import (
	"testing"

	"github.com/powerman/go-monolith-example/internal/apiauth"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/api"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/ms/example/internal/config"
	"github.com/powerman/go-monolith-example/ms/example/internal/dal"
	"github.com/powerman/gotest/testinit"
)

func TestMain(m *testing.M) { testinit.Main(m) }

const (
	serialMain = iota
	serialIntegration
)

func init() { testinit.Setup(serialMain, setupMain) }

func setupMain() {
	def.Init()
	dal.InitMetrics(reg)
	app.InitMetrics(reg)
	api.InitMetrics(reg)
	cfg = config.MustGetTest()
}

var (
	ctx        = def.NewContext(app.ServiceName)
	tokenAdmin = apiauth.AccessToken("admin")
	tokenUser  = apiauth.AccessToken("user")
	authAdmin  = dom.Auth{
		UserID: 1,
		Admin:  true,
	}
	authUser = dom.Auth{
		UserID: 2,
		Admin:  false,
	}
)
