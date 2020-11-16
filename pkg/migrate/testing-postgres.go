package migrate

import (
	"runtime"
	"strings"

	"github.com/powerman/check"
	goosepkg "github.com/powerman/goose/v2"
	"github.com/powerman/pqx"

	"github.com/powerman/go-monolith-example/pkg/def"
)

// PostgresUpDownTest creates temporary database, test given migrations, and removes
// temporary database.
func PostgresUpDownTest(t *check.C, ctx Ctx, goose *goosepkg.Instance, dir string, cfg *def.PostgresConfig) {
	pc, _, _, _ := runtime.Caller(1)
	suffix := runtime.FuncForPC(pc).Name()
	suffix = suffix[:strings.LastIndex(suffix, ".")]

	_, cleanup, err := pqx.EnsureTempDB(tLogger(*t), suffix, cfg.Config)
	t.Must(t.Nil(err))
	defer cleanup()

	connector := &Postgres{PostgresConfig: cfg.Clone()}
	connector.PostgresConfig.DBName += "_" + suffix
	upDownTest(t, ctx, goose, dir, connector)
}
