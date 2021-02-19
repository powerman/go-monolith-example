package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/go-sql-driver/mysql"
	goosepkg "github.com/powerman/goose/v2"
	"github.com/powerman/must"

	// Driver.
	_ "github.com/powerman/narada4d/protocol/goose-mysql"
	"github.com/powerman/narada4d/schemaver"
	"github.com/powerman/structlog"
)

var reTCP = regexp.MustCompile(`(^|@)tcp[(]([^)]*)[)]`)

// MySQL implements Connector interface.
type MySQL struct {
	*mysql.Config
}

// Connect to MySQL. Will create database and initialize schemaver if needed.
func (c *MySQL) Connect(ctx Ctx, goose *goosepkg.Instance) (_ *sql.DB, _ *schemaver.SchemaVer, err error) {
	log := structlog.FromContext(ctx, nil)

	cfg := c.Clone()
	cfg.MaxAllowedPacket = 0
	cfg.MultiStatements = true // https://github.com/pressly/goose/issues/190

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, nil, fmt.Errorf("sql.Open: %w", err)
	}
	defer func() {
		if err != nil {
			log.WarnIfFail(db.Close)
		}
	}()

	if cfg.Timeout != 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}
	err = db.PingContext(ctx)
	if err2 := new(mysql.MySQLError); errors.As(err, &err2) && err2.Number == 1049 {
		cfgNoDB := cfg.Clone()
		cfgNoDB.DBName = ""
		db2, err := sql.Open("mysql", cfgNoDB.FormatDSN())
		if err != nil {
			return nil, nil, fmt.Errorf("sql.Open: %w", err)
		}
		_, err = db2.ExecContext(ctx, fmt.Sprintf(
			"CREATE DATABASE IF NOT EXISTS `%s` COLLATE %s", cfg.DBName, cfg.Collation))
		log.WarnIfFail(db2.Close)
		if err != nil {
			return nil, nil, fmt.Errorf("create database %q: %w", cfg.DBName, err)
		}
	}
	for err != nil {
		nextErr := db.PingContext(ctx)
		if errors.Is(nextErr, context.DeadlineExceeded) || errors.Is(nextErr, context.Canceled) {
			return nil, nil, fmt.Errorf("db.Ping: %w", err)
		}
		err = nextErr
	}

	gooseMu.Lock()
	defer gooseMu.Unlock()
	must.NoErr(goose.SetDialect("mysql"))
	_, _ = goose.EnsureDBVersion(db) // Race on CREATE TABLE, so allowed to fail.

	ver, err := schemaver.NewAt("goose-mysql://" + reTCP.ReplaceAllString(cfg.FormatDSN(), "$1$2"))
	if err != nil {
		return nil, nil, err
	}

	return db, ver, nil
}
