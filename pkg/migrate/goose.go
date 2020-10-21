// Package migrate manage DB migrations.
package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-sql-driver/mysql"
	goosepkg "github.com/powerman/goose/v2"
	"github.com/powerman/must"

	// Driver.
	_ "github.com/powerman/narada4d/protocol/goose-mysql"
	"github.com/powerman/narada4d/schemaver"
	"github.com/powerman/structlog"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

var errSelfCheck = errors.New("unexpected db schema version")

//nolint:gochecknoglobals // Regexp.
var reTCP = regexp.MustCompile(`(^|@)tcp[(]([^)]*)[)]`)

func connect(ctx Ctx, goose *goosepkg.Instance, cfg *mysql.Config) (db *sql.DB, ver *schemaver.SchemaVer, err error) {
	log := structlog.FromContext(ctx, nil)

	cfg = cfg.Clone()
	cfg.MaxAllowedPacket = 0
	cfg.MultiStatements = true // https://github.com/pressly/goose/issues/190

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, nil, fmt.Errorf("sql.Open: %w", err)
	}
	defer func(dbClose func() error) {
		if err != nil {
			log.WarnIfFail(dbClose)
		}
	}(db.Close)

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

	must.NoErr(goose.SetDialect("mysql"))
	_, _ = goose.EnsureDBVersion(db) // Race on CREATE TABLE, so allowed to fail.

	ver, err = schemaver.NewAt("goose-mysql://" + reTCP.ReplaceAllString(cfg.FormatDSN(), "$1$2"))
	if err != nil {
		return nil, nil, err
	}

	return db, ver, nil
}

// UpTo migrates up to a specific version.
//
// Unlike goose.UpTo it will return error is current version doesn't match
// requested one after migration.
func UpTo(ctx Ctx, goose *goosepkg.Instance, dir string, version int64, cfg *mysql.Config) (*schemaver.SchemaVer, error) {
	log := structlog.FromContext(ctx, nil)

	db, ver, err := connect(ctx, goose, cfg)
	if err != nil {
		return nil, err
	}
	defer log.WarnIfFail(db.Close)

	_ = ver.ExclusiveLock()
	defer ver.Unlock()

	err = goose.UpTo(db, dir, version)
	if err != nil {
		return nil, fmt.Errorf("goose.UpTo %d: %w", version, err)
	}
	if v, _ := goose.GetDBVersion(db); v != version {
		return nil, fmt.Errorf("%w: %d (should be %d)", errSelfCheck, v, version)
	}

	return ver, nil
}

// Run executes goose command. It also enforce "fix" after "create".
func Run(ctx Ctx, goose *goosepkg.Instance, dir string, command string, cfg *mysql.Config) error {
	log := structlog.FromContext(ctx, nil)

	db, ver, err := connect(ctx, goose, cfg)
	if err != nil {
		return err
	}
	defer log.WarnIfFail(db.Close)
	defer log.WarnIfFail(ver.Close)

	_ = ver.ExclusiveLock()
	defer ver.Unlock()

	cmdArgs := strings.Fields(command)
	cmd, args := cmdArgs[0], cmdArgs[1:]
	err = goose.Run(cmd, db, dir, args...)
	if err == nil && cmd == "create" {
		err = goose.Run("fix", db, dir)
	}
	if err != nil {
		return fmt.Errorf("goose.Run %q: %w", command, err)
	}
	return nil
}
