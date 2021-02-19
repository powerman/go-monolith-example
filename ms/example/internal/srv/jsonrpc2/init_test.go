package jsonrpc2_test

import (
	"net/http/httptest"
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
	"github.com/powerman/go-monolith-example/pkg/jsonrpc2x"
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
	tokenEmpty = ""
	tokenAdmin = "admin"
	tokenUser  = "user"
	authAdmin  = dom.Auth{
		UserName: dom.NewUserName("1"),
		Admin:    true,
	}
	authUser = dom.Auth{
		UserName: dom.NewUserName("2"),
		Admin:    false,
	}
	userIDBad = dom.UserName{Name: dom.NewName("guests", "0")}
)

func testNew(t *check.C) (*jsonrpc2x.Client, string, *app.MockAppl) {
	ctrl := gomock.NewController(t)

	mockAppl := app.NewMockAppl(ctrl)
	mockAuthn := apix.NewMockAuthn(ctrl)
	srv := jsonrpc2.NewServer(mockAppl, mockAuthn, jsonrpc2.Config{
		Pattern:   "/",
		StrictErr: true,
	})

	mockAuthn.EXPECT().Authenticate(gomock.Any(), apix.AccessToken(tokenAdmin)).Return(authAdmin, nil).AnyTimes()
	mockAuthn.EXPECT().Authenticate(gomock.Any(), apix.AccessToken(tokenUser)).Return(authUser, nil).AnyTimes()
	mockAuthn.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(dom.Auth{}, apix.ErrAccessTokenInvalid).AnyTimes()

	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)
	return jsonrpc2x.NewHTTPClient(ts.URL), ts.URL, mockAppl
}
