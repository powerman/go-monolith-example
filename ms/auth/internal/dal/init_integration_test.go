// +build integration

package dal_test

import (
	"errors"
	"runtime"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/powerman/check"
	"github.com/powerman/pqx"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/ms/auth/internal/config"
	"github.com/powerman/go-monolith-example/ms/auth/internal/dal"
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
	const maxIdentLen = 63
	if maxLen := maxIdentLen - len(cfg.Postgres.DBName) - 1; len(suffix) > maxLen {
		suffix = suffix[len(suffix)-maxLen:]
	}

	dbCfg := cfg.Postgres.Clone()
	_, cleanupDB, err := pqx.EnsureTempDB(tLogger(*t), suffix, dbCfg.Config)
	t.Must(t.Nil(err))
	t.Cleanup(cleanupDB)
	tempDBCfg := dbCfg.Clone()
	tempDBCfg.DBName += "_" + suffix
	r, err := dal.New(ctx, cfg.GoosePostgresDir, tempDBCfg)
	t.Must(t.Nil(err))
	t.Cleanup(r.Close)

	return r
}

func matchErr(t *check.C, err, wantErr error) {
	t.Helper()
	if pqErr := new(*pq.Error); errors.As(err, pqErr) && wantErr != nil {
		t.Match(err, wantErr.Error())
	} else {
		t.Err(err, wantErr)
	}
}
