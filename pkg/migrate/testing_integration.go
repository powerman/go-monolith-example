// +build integration

package migrate

import (
	"runtime"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/check"
	goosepkg "github.com/powerman/goose/v2"
	"github.com/powerman/mysqlx"
)

type tLogger check.C

func (t tLogger) Print(args ...interface{}) { t.Log(args...) }

// UpDownTest creates temporary database, test given migrations, and removes
// temporary database.
func UpDownTest(t *check.C, ctx Ctx, goose *goosepkg.Instance, dir string, cfg *mysql.Config) {
	pc, _, _, _ := runtime.Caller(1)
	suffix := runtime.FuncForPC(pc).Name()
	suffix = suffix[:strings.LastIndex(suffix, ".")]

	cfg, cleanup, err := mysqlx.EnsureTempDB(tLogger(*t), suffix, cfg)
	t.Must(t.Nil(err))
	defer cleanup()

	db, _, err := connect(ctx, goose, cfg)
	t.Must(t.Nil(err))
	defer db.Close()

	t.Must(t.Nil(Run(ctx, goose, dir, "up", cfg)))
	for v, _ := goose.GetDBVersion(db); v > 0; v, _ = goose.GetDBVersion(db) {
		err := Run(ctx, goose, dir, "down", cfg)
		if err != nil && t.Contains(err.Error(), ErrDownNotSupported.Error()) {
			t.Logf("downgrade from version %d is not supported", v)
			t.Nil(Run(ctx, goose, dir, "up", cfg))
			return
		}
		t.Must(t.Nil(err))
		v2, err := goose.GetDBVersion(db)
		t.Nil(err)
		t.Less(v2, v)
	}
	v, err := goose.GetDBVersion(db)
	t.Nil(err)
	t.Zero(v)
	t.Nil(Run(ctx, goose, dir, "up", cfg))
}
