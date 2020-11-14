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
	"github.com/powerman/sqlxx"
	"github.com/powerman/structlog"

	"github.com/powerman/go-monolith-example/pkg/migrate"
)

// MySQLMaxKeySize for indexed MySQL utf8mb4 CHAR/VARCHAR column.
const MySQLMaxKeySize = 191

// MySQLDuplicateEntry returns true if err is mysql error "Duplicate entry…".
func MySQLDuplicateEntry(err error) bool {
	const duplicateEntry = 1062
	if errMySQL := new(mysql.MySQLError); errors.As(err, &errMySQL) {
		return errMySQL.Number == duplicateEntry
	}
	return false
}

// MySQLConfig contains repo configuration.
type MySQLConfig struct {
	MySQL         *mysql.Config
	GooseMySQLDir string
	SchemaVersion int64
	Metric        Metrics
	ReturnErrs    []error // List of app.Err… returned by DAL methods.
}

// NewMySQL creates and returns new Repo.
// It will also run required DB migrations and connects to DB.
func NewMySQL(ctx Ctx, goose *goosepkg.Instance, cfg MySQLConfig) (*Repo, error) {
	log := structlog.FromContext(ctx, nil)

	connector := &migrate.MySQL{Config: cfg.MySQL}
	schemaVer, err := migrate.UpTo(ctx, goose, cfg.GooseMySQLDir, cfg.SchemaVersion, connector)
	if err != nil {
		return nil, fmt.Errorf("migration: %w", err)
	}
	defer func() {
		if err != nil {
			log.WarnIfFail(schemaVer.Close)
		}
	}()

	db, err := sql.Open("mysql", cfg.MySQL.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	defer func() {
		if err != nil {
			log.WarnIfFail(db.Close)
		}
	}()

	if cfg.MySQL.Timeout != 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, cfg.MySQL.Timeout)
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
		DB:            sqlxx.NewDB(sqlx.NewDb(db, "mysql")),
		SchemaVer:     schemaVer,
		schemaVersion: strconv.Itoa(int(cfg.SchemaVersion)),
		returnErrs:    cfg.ReturnErrs,
		metric:        cfg.Metric,
		log:           log,
		serialize:     mysqlSerialize,
	}
	return r, nil
}

func mysqlSerialize(doTx func() error) error {
	// TODO Implement auto-retries if transaction fails because of
	// serialization-related error.
	return doTx()
}
