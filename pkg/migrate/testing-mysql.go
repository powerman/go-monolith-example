package migrate

import (
	"runtime"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/check"
	goosepkg "github.com/powerman/goose/v2"
	"github.com/powerman/mysqlx"
)

// MySQLUpDownTest creates temporary database, test given migrations, and removes
// temporary database.
func MySQLUpDownTest(t *check.C, ctx Ctx, goose *goosepkg.Instance, dir string, cfg *mysql.Config) {
	pc, _, _, _ := runtime.Caller(1)
	suffix := runtime.FuncForPC(pc).Name()
	suffix = suffix[:strings.LastIndex(suffix, ".")]

	cfg, cleanup, err := mysqlx.EnsureTempDB(tLogger(*t), suffix, cfg)
	t.Must(t.Nil(err))
	defer cleanup()

	connector := &MySQL{Config: cfg}
	upDownTest(t, ctx, goose, dir, connector)
}
