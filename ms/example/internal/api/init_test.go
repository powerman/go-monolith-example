package api

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/proto"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	reg := prometheus.NewPedanticRegistry()
	def.Init()
	app.InitMetrics(reg)
	InitMetrics(reg)
	check.TestMain(m)
}

// Const shared by tests. Recommended naming scheme: <dataType><Variant>.
var (
	tokenEmpty = proto.AccessToken("")
	tokenAdmin = proto.AccessToken("admin")
	tokenUser  = proto.AccessToken("user")
	authAdmin  = dom.Auth{
		UserID: 1,
		Admin:  true,
	}
	authUser = dom.Auth{
		UserID: 2,
		Admin:  false,
	}
	userIDBad = dom.UserID(0)
)

func testNew(t *check.C) (*API, *app.MockAppl) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockApp := app.NewMockAppl(ctrl)
	api := New(mockApp)
	return api, mockApp
}
