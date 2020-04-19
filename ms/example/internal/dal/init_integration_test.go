// +build integration

package dal

import (
	"github.com/powerman/gotest/testinit"
	"github.com/powerman/mysqlx"
	"github.com/powerman/structlog"
)

var r *Repo

func init() { testinit.Setup(serialIntegration, setupIntegration) }

func setupIntegration() {
	const dir = "../migrations"
	log := structlog.FromContext(ctx, nil)

	cfg, cleanup, err := mysqlx.EnsureTempDB(log, "", cfg.MySQLConfig)
	if err == nil {
		testinit.Teardown(cleanup)
		r, err = New(ctx, dir, cfg)
	}
	if err != nil {
		testinit.Fatal(err)
	}
	testinit.Teardown(r.Close)
}
