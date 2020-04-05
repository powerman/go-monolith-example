// Package dal implements Data Access Layer.
package dal

import (
	"context"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/go-monolith-example/internal/repo"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/ms/example/internal/migrations"
)

const (
	schemaVersion  = 3
	dbMaxOpenConns = 0 // unlimited
	dbMaxIdleConns = 5 // a bit more than default (2)
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Repo provides access to storage.
type Repo struct {
	*repo.Repo
}

// New creates and returns new Repo.
// It will also run required DB migrations and connects to DB.
func New(ctx Ctx, dir string, cfg *mysql.Config) (_ *Repo, err error) {
	var returnErrs = []error{ // List of app.Err… returned by Repo methods.
		app.ErrNotFound,
	}

	r := &Repo{}
	r.Repo, err = repo.New(ctx, migrations.Goose(), dir, schemaVersion, returnErrs, metric, cfg)
	if err != nil {
		return nil, err
	}
	r.DB.SetMaxOpenConns(dbMaxOpenConns)
	r.DB.SetMaxIdleConns(dbMaxIdleConns)
	r.SchemaVer.HoldSharedLock(ctx, time.Second)
	return r, nil
}
