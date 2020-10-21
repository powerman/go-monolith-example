// Package repo provide helpers for Data Access Layer.
package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	goosepkg "github.com/powerman/goose/v2"
	"github.com/powerman/narada4d/schemaver"
	"github.com/powerman/sqlxx"
	"github.com/powerman/structlog"

	"github.com/powerman/go-monolith-example/pkg/migrate"
	"github.com/powerman/go-monolith-example/pkg/reflectx"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// MaxKeySize for indexed MySQL utf8mb4 CHAR/VARCHAR column.
const MaxKeySize = 191

// Errors.
var (
	ErrSchemaVer = errors.New("unsupported DB schema version")
)

// DuplicateEntry returns true if err is mysql error "Duplicate entry…".
func DuplicateEntry(err error) bool {
	const duplicateEntry = 1062
	if errMySQL := new(mysql.MySQLError); errors.As(err, &errMySQL) {
		return errMySQL.Number == duplicateEntry
	}
	return false
}

// Config contains repo configuration.
type Config struct {
	MySQL         *mysql.Config
	GooseDir      string
	SchemaVersion int64
	Metric        Metrics
	ReturnErrs    []error // List of app.Err… returned by DAL methods.
}

// Repo provides access to storage.
type Repo struct {
	DB            *sqlxx.DB
	SchemaVer     *schemaver.SchemaVer
	schemaVersion string
	returnErrs    []error
	metric        Metrics
	log           *structlog.Logger
}

// New creates and returns new Repo.
// It will also run required DB migrations and connects to DB.
func New(ctx Ctx, goose *goosepkg.Instance, cfg Config) (*Repo, error) {
	log := structlog.FromContext(ctx, nil)

	schemaVer, err := migrate.UpTo(ctx, goose, cfg.GooseDir, cfg.SchemaVersion, cfg.MySQL)
	if err != nil {
		return nil, fmt.Errorf("migration: %w", err)
	}

	db, err := sql.Open("mysql", cfg.MySQL.FormatDSN())
	if err != nil {
		log.WarnIfFail(schemaVer.Close)
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if cfg.MySQL.Timeout != 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, cfg.MySQL.Timeout)
		defer cancel()
	}
	err = db.PingContext(ctx)
	for err != nil {
		nextErr := db.PingContext(ctx)
		if errors.Is(nextErr, context.DeadlineExceeded) || errors.Is(nextErr, context.Canceled) {
			log.WarnIfFail(db.Close)
			log.WarnIfFail(schemaVer.Close)
			return nil, fmt.Errorf("db.Ping: %w", err)
		}
		err = nextErr
	}

	r := &Repo{
		DB:            sqlxx.NewDB(sqlx.NewDb(db, "mysql")),
		SchemaVer:     schemaVer,
		schemaVersion: strconv.Itoa(int(cfg.SchemaVersion)),
		returnErrs:    cfg.ReturnErrs,
		metric:        cfg.Metric,
		log:           log,
	}
	return r, nil
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

func (r *Repo) schemaLock(f func() error) error {
	ver := r.SchemaVer.SharedLock()
	defer r.SchemaVer.Unlock()
	if ver != r.schemaVersion {
		return fmt.Errorf("schema version %s, need %s: %w", ver, r.schemaVersion, ErrSchemaVer)
	}
	return f()
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
	})))
}

// Tx provides DAL method wrapper with:
// - converting sqlx errors which are actually bugs into panics,
// - ensure valid schema version while accessing DB,
// - general metrics for DAL methods,
// - wrapping errors with DAL method name,
// - transaction.
func (r *Repo) Tx(ctx Ctx, opts *sql.TxOptions, f func(*sqlxx.Tx) error) (err error) {
	methodName := reflectx.CallerMethodName(1)
	return r.strict(r.schemaLock(r.metric.instrument(methodName, func() error {
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
	})))
}
