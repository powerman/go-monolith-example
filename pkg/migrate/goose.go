// Package migrate manage DB migrations.
package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	goosepkg "github.com/powerman/goose/v2"
	"github.com/powerman/narada4d/schemaver"
	"github.com/powerman/structlog"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

var errSelfCheck = errors.New("unexpected db schema version")

// Tests often runs in parallel using same goose instance and may trigger
// -race detector on SetDialect. So, use this mutex to work around.
//nolint:gochecknoglobals // By design.
var gooseMu sync.Mutex

// Connector provides a way to connect to any database with schemaver.
type Connector interface {
	Connect(Ctx, *goosepkg.Instance) (*sql.DB, *schemaver.SchemaVer, error)
}

// UpTo migrates up to a specific version.
//
// Unlike goose.UpTo it will return error is current version doesn't match
// requested one after migration.
func UpTo(ctx Ctx, goose *goosepkg.Instance, dir string, version int64, c Connector) (*schemaver.SchemaVer, error) {
	log := structlog.FromContext(ctx, nil)

	db, ver, err := c.Connect(ctx, goose)
	if err != nil {
		return nil, err
	}
	defer log.WarnIfFail(db.Close)

	_ = ver.ExclusiveLock()
	defer ver.Unlock()

	gooseMu.Lock()
	defer gooseMu.Unlock()
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
func Run(ctx Ctx, goose *goosepkg.Instance, dir string, command string, c Connector) error {
	log := structlog.FromContext(ctx, nil)

	db, ver, err := c.Connect(ctx, goose)
	if err != nil {
		return err
	}
	defer log.WarnIfFail(db.Close)
	defer log.WarnIfFail(ver.Close)

	_ = ver.ExclusiveLock()
	defer ver.Unlock()

	gooseMu.Lock()
	defer gooseMu.Unlock()
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
