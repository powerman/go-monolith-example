package app_test

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/smartystreets/goconvey/convey"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/pkg/def"
)

func TestMain(m *testing.M) {
	def.Init()
	reg := prometheus.NewPedanticRegistry()
	app.InitMetrics(reg)
	check.TestMain(m)
}

type Ctx = context.Context

// Const shared by tests. Recommended naming scheme: <dataType><Variant>.
var (
	ctx       = def.NewContext(app.ServiceName)
	userIDBad = dom.UserID(666)
	authAdmin = dom.Auth{
		UserID: 1,
		Admin:  true,
	}
	authUser = dom.Auth{
		UserID: 2,
		Admin:  false,
	}
)

func testNew(t *check.C) (func(), *app.App, *app.MockRepo) {
	ctrl := gomock.NewController(t)

	mockRepo := app.NewMockRepo(ctrl)
	a := app.New(mockRepo, app.Config{})
	return ctrl.Finish, a, mockRepo
}
