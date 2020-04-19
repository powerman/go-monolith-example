package api

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/go-monolith-example/internal/apiauth"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
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
	tokenEmpty = apiauth.AccessToken("")
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
	userIDBad = dom.UserID(0)
)

func testNew(t *check.C) (*API, *app.MockAppl) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockApp := app.NewMockAppl(ctrl)
	mockAuthn := apiauth.NewMockAuthenticator(ctrl)
	api := New(mockApp, mockAuthn)
	api.strictErr = true

	mockAuthn.EXPECT().Authenticate(tokenAdmin).Return(authAdmin, nil).AnyTimes()
	mockAuthn.EXPECT().Authenticate(tokenUser).Return(authUser, nil).AnyTimes()
	mockAuthn.EXPECT().Authenticate(gomock.Any()).Return(dom.Auth{}, apiauth.ErrAccessTokenInvalid).AnyTimes()

	return api, mockApp
}
