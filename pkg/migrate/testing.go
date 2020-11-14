package migrate

import (
	"github.com/powerman/check"
	goosepkg "github.com/powerman/goose/v2"
)

type tLogger check.C

func (t tLogger) Print(args ...interface{}) { t.Log(args...) }

func upDownTest(t *check.C, ctx Ctx, goose *goosepkg.Instance, dir string, connector Connector) {
	db, ver, err := connector.Connect(ctx, goose)
	t.Must(t.Nil(err))
	defer func() {
		t.Nil(db.Close())
		t.Nil(ver.Close())
	}()

	t.Must(t.Nil(Run(ctx, goose, dir, "up", connector)))
	for v, _ := goose.GetDBVersion(db); v > 0; v, _ = goose.GetDBVersion(db) {
		err := Run(ctx, goose, dir, "down", connector)
		if err != nil && t.Contains(err.Error(), ErrDownNotSupported.Error()) {
			t.Logf("downgrade from version %d is not supported", v)
			t.Nil(Run(ctx, goose, dir, "up", connector))
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
	t.Nil(Run(ctx, goose, dir, "up", connector))
}
