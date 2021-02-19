// +build integration

package dal_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/mysqlx"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/ms/example/internal/config"
	"github.com/powerman/go-monolith-example/ms/example/internal/dal"
	"github.com/powerman/go-monolith-example/pkg/def"
)

func TestMain(m *testing.M) {
	def.Init()
	reg := prometheus.NewPedanticRegistry()
	app.InitMetrics(reg)
	dal.InitMetrics(reg)
	cfg = config.MustGetServeTest()
	check.TestMain(m)
}

type tLogger check.C

func (t tLogger) Print(args ...interface{}) { t.Log(args...) }

var (
	ctx = def.NewContext(app.ServiceName)
	cfg *config.ServeConfig
)

func newTestRepo(t *check.C) *dal.Repo {
	t.Helper()

	pc, _, _, _ := runtime.Caller(1)
	suffix := runtime.FuncForPC(pc).Name()
	suffix = suffix[:strings.LastIndex(suffix, ".")]
	suffix += "_" + t.Name()

	tempDBCfg, cleanupDB, err := mysqlx.EnsureTempDB(tLogger(*t), suffix, cfg.MySQL)
	t.Must(t.Nil(err))
	t.Cleanup(cleanupDB)
	r, err := dal.New(ctx, cfg.GooseMySQLDir, tempDBCfg)
	t.Must(t.Nil(err))
	t.Cleanup(r.Close)

	return r
}
