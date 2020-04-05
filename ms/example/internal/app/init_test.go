package app

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	reg := prometheus.NewPedanticRegistry()
	def.Init()
	InitMetrics(reg)
	check.TestMain(m)
}

// Const shared by tests. Recommended naming scheme: <dataType><Variant>.
var (
	ctx       = def.NewContext(ServiceName)
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

func testNew(t *check.C) (*App, *MockRepo) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := NewMockRepo(ctrl)
	a := New(mockRepo, Config{})
	return a, mockRepo
}
