package jsonrpc2_test

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/smartystreets/goconvey/convey"

	"github.com/powerman/go-monolith-example/internal/apix"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/ms/example/internal/srv/jsonrpc2"
	"github.com/powerman/go-monolith-example/pkg/def"
)

func TestMain(m *testing.M) {
	reg := prometheus.NewPedanticRegistry()
	def.Init()
	app.InitMetrics(reg)
	jsonrpc2.InitMetrics(reg)
	check.TestMain(m)
}

// Const shared by tests. Recommended naming scheme: <dataType><Variant>.
var (
	tokenEmpty = apix.AccessToken("")
	tokenAdmin = apix.AccessToken("admin")
	tokenUser  = apix.AccessToken("user")
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

func testNew(t *check.C) (*jsonrpc2.Server, *app.MockAppl) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockApp := app.NewMockAppl(ctrl)
	mockAuthn := apix.NewMockAuthn(ctrl)
	srv := jsonrpc2.New(mockApp, mockAuthn, jsonrpc2.Config{
		StrictErr: true,
	})

	mockAuthn.EXPECT().Authenticate(tokenAdmin).Return(authAdmin, nil).AnyTimes()
	mockAuthn.EXPECT().Authenticate(tokenUser).Return(authUser, nil).AnyTimes()
	mockAuthn.EXPECT().Authenticate(gomock.Any()).Return(dom.Auth{}, apix.ErrAccessTokenInvalid).AnyTimes()

	return srv, mockApp
}