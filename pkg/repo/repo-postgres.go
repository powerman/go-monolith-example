package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	goosepkg "github.com/powerman/goose/v2"
	"github.com/powerman/pqx"
	"github.com/powerman/sqlxx"
	"github.com/powerman/structlog"

	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/migrate"
)

// Error names.
const (
	PostgresUniqueViolation     = "unique_violation"
	PostgresForeignKeyViolation = "foreign_key_violation"
)

// PostgresErrName returns true if err is PostgreSQL error with given name.
func PostgresErrName(err error, name string) bool {
	pqErr := new(pq.Error)
	return errors.As(err, &pqErr) && pqErr.Code.Name() == name
}

// PostgresConfig contains repo configuration.
type PostgresConfig struct {
	Postgres         *def.PostgresConfig
	GoosePostgresDir string
	SchemaVersion    int64
	Metric           Metrics
	ReturnErrs       []error // List of app.Err… returned by DAL methods.
}

// NewPostgres creates and returns new Repo.
// It will also run required DB migrations and connects to DB.
func NewPostgres(ctx Ctx, goose *goosepkg.Instance, cfg PostgresConfig) (_ *Repo, err error) {
	log := structlog.FromContext(ctx, nil)

	connector := &migrate.Postgres{PostgresConfig: cfg.Postgres}
	schemaVer, err := migrate.UpTo(ctx, goose, cfg.GoosePostgresDir, cfg.SchemaVersion, connector)
	if err != nil {
		return nil, fmt.Errorf("migration: %w", err)
	}
	defer func() {
		if err != nil {
			log.WarnIfFail(schemaVer.Close)
		}
	}()

	err = cfg.Postgres.UpdateConnectTimeout(ctx)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", cfg.Postgres.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	defer func() {
		if err != nil {
			log.WarnIfFail(db.Close)
		}
	}()

	if cfg.Postgres.ConnectTimeout != 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, cfg.Postgres.ConnectTimeout)
		defer cancel()
	}
	err = db.PingContext(ctx)
	for err != nil {
		nextErr := db.PingContext(ctx)
		if errors.Is(nextErr, context.DeadlineExceeded) || errors.Is(nextErr, context.Canceled) {
			return nil, fmt.Errorf("db.Ping: %w", err)
		}
		err = nextErr
	}

	r := &Repo{
		DB:            sqlxx.NewDB(sqlx.NewDb(db, "postgres")),
		SchemaVer:     schemaVer,
		schemaVersion: strconv.Itoa(int(cfg.SchemaVersion)),
		returnErrs:    cfg.ReturnErrs,
		metric:        cfg.Metric,
		log:           log,
		serialize:     pqx.Serialize,
	}
	return r, nil
}
