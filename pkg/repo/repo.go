// Package repo provide helpers for Data Access Layer.
package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/powerman/narada4d/schemaver"
	"github.com/powerman/sqlxx"
	"github.com/powerman/structlog"

	"github.com/powerman/go-monolith-example/pkg/reflectx"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Errors.
var (
	ErrSchemaVer = errors.New("unsupported DB schema version")
)

// Repo provides access to storage.
type Repo struct {
	DB            *sqlxx.DB
	SchemaVer     *schemaver.SchemaVer
	schemaVersion string
	returnErrs    []error
	metric        Metrics
	log           *structlog.Logger
	serialize     func(doTx func() error) error
}

// Close closes connection to DB.
func (r *Repo) Close() {
	r.log.WarnIfFail(r.DB.Close)
	r.log.WarnIfFail(r.SchemaVer.Close)
}

// Turn sqlx errors like `missing destination …` into panics
// https://github.com/jmoiron/sqlx/issues/529. As we can't distinguish
// between sqlx and other errors except driver ones, let's hope filtering
// driver errors is enough and there are no other non-driver regular errors.
func (r *Repo) strict(err error) error {
	switch {
	case err == nil:
	case errors.As(err, new(*mysql.MySQLError)):
	case errors.As(err, new(*pq.Error)):
	case errors.Is(err, ErrSchemaVer):
	case errors.Is(err, sql.ErrNoRows):
	case errors.Is(err, context.Canceled):
	case errors.Is(err, context.DeadlineExceeded):
	default:
		for i := range r.returnErrs {
			if errors.Is(err, r.returnErrs[i]) {
				return err
			}
		}
		panic(err)
	}
	return err
}

func (r *Repo) schemaLock(f func() error) func() error {
	return func() error {
		ver := r.SchemaVer.SharedLock()
		defer r.SchemaVer.Unlock()
		if ver != r.schemaVersion {
			return fmt.Errorf("schema version %s, need %s: %w", ver, r.schemaVersion, ErrSchemaVer)
		}
		return f()
	}
}

// NoTx provides DAL method wrapper with:
// - converting sqlx errors which are actually bugs into panics,
// - ensure valid schema version while accessing DB,
// - general metrics for DAL methods,
// - wrapping errors with DAL method name.
func (r *Repo) NoTx(f func() error) (err error) {
	methodName := reflectx.CallerMethodName(1)
	return r.strict(r.schemaLock(r.metric.instrument(methodName, func() error {
		err := f()
		if err != nil {
			err = fmt.Errorf("%s: %w", methodName, err)
		}
		return err
	}))())
}

// Tx provides DAL method wrapper with:
// - converting sqlx errors which are actually bugs into panics,
// - ensure valid schema version while accessing DB,
// - general metrics for DAL methods,
// - wrapping errors with DAL method name,
// - transaction.
func (r *Repo) Tx(ctx Ctx, opts *sql.TxOptions, f func(*sqlxx.Tx) error) (err error) {
	methodName := reflectx.CallerMethodName(1)
	return r.strict(r.serialize(r.schemaLock(r.metric.instrument(methodName, func() error {
		tx, err := r.DB.BeginTxx(ctx, opts)
		if err == nil { //nolint:nestif // No idea how to simplify.
			defer func() {
				if err := recover(); err != nil {
					if err := tx.Rollback(); err != nil {
						log := structlog.FromContext(ctx, nil)
						log.Warn("failed to tx.Rollback", "method", methodName, "err", err)
					}
					panic(err)
				}
			}()
			err = f(tx)
			if err == nil {
				err = tx.Commit()
			} else if err := tx.Rollback(); err != nil {
				log := structlog.FromContext(ctx, nil)
				log.Warn("failed to tx.Rollback", "method", methodName, "err", err)
			}
		}
		if err != nil {
			err = fmt.Errorf("%s: %w", methodName, err)
		}
		return err
	}))))
}
