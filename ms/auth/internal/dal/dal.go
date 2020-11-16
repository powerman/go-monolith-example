// Package dal implements Data Access Layer using PostgreSQL DB.
package dal

import (
	"context"
	"time"

	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/ms/auth/internal/migrations"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/repo"
)

const (
	schemaVersion  = 4
	dbMaxOpenConns = 100 / 10 // Use up to 1/10 of server's max_connections.
	dbMaxIdleConns = 5        // A bit more than default (2).
)

type Ctx = context.Context

type Repo struct {
	*repo.Repo
}

// New creates and returns new Repo.
// It will also run required DB migrations and connects to DB.
func New(ctx Ctx, dir string, cfg *def.PostgresConfig) (_ *Repo, err error) {
	returnErrs := []error{ // List of app.Err… returned by Repo methods.
		app.ErrAlreadyExist,
		app.ErrNotFound,
	}

	r := &Repo{}
	r.Repo, err = repo.NewPostgres(ctx, migrations.Goose(), repo.PostgresConfig{
		Postgres:         cfg,
		GoosePostgresDir: dir,
		SchemaVersion:    schemaVersion,
		Metric:           metric,
		ReturnErrs:       returnErrs,
	})
	if err != nil {
		return nil, err
	}
	r.DB.SetMaxOpenConns(dbMaxOpenConns)
	r.DB.SetMaxIdleConns(dbMaxIdleConns)
	r.SchemaVer.HoldSharedLock(ctx, time.Second)
	return r, nil
}
