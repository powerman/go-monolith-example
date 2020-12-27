package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	goosepkg "github.com/powerman/goose/v2"
	"github.com/powerman/must"

	// Driver.
	_ "github.com/powerman/narada4d/protocol/goose-postgres"
	"github.com/powerman/narada4d/schemaver"
	"github.com/powerman/structlog"

	"github.com/powerman/go-monolith-example/pkg/def"
)

const (
	postgresStatementTimeout                = time.Second
	postgresLockTimeout                     = time.Second
	postgresIdleInTransactionSessionTimeout = time.Second
)

// Postgres implements Connector interface.
type Postgres struct {
	*def.PostgresConfig
}

// Connect to PostgreSQL. Will initialize schemaver if needed.
func (c *Postgres) Connect(ctx Ctx, goose *goosepkg.Instance) (_ *sql.DB, _ *schemaver.SchemaVer, err error) {
	log := structlog.FromContext(ctx, nil)

	cfg := c.PostgresConfig.Clone()
	cfg.DefaultTransactionIsolation = sql.LevelDefault
	cfg.StatementTimeout = postgresStatementTimeout
	cfg.LockTimeout = postgresLockTimeout
	cfg.IdleInTransactionSessionTimeout = postgresIdleInTransactionSessionTimeout
	err = cfg.UpdateConnectTimeout(ctx)
	if err != nil {
		return nil, nil, err
	}

	db, err := sql.Open("postgres", cfg.FormatDSN())
	if err != nil {
		return nil, nil, fmt.Errorf("sql.Open: %w", err)
	}
	defer func() {
		if err != nil {
			log.WarnIfFail(db.Close)
		}
	}()

	if cfg.ConnectTimeout != 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, cfg.ConnectTimeout)
		defer cancel()
	}
	err = db.PingContext(ctx)
	for err != nil {
		nextErr := db.PingContext(ctx)
		if errors.Is(nextErr, context.DeadlineExceeded) || errors.Is(nextErr, context.Canceled) {
			return nil, nil, fmt.Errorf("db.Ping: %w", err)
		}
		err = nextErr
	}

	gooseMu.Lock()
	defer gooseMu.Unlock()
	must.NoErr(goose.SetDialect("postgres"))
	_, _ = goose.EnsureDBVersion(db) // Race on CREATE TABLE, so allowed to fail.

	ver, err := schemaver.NewAt("goose-" + cfg.FormatURL())
	if err != nil {
		return nil, nil, err
	}

	return db, ver, nil
}
