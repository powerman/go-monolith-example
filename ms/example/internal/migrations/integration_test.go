// +build integration

package migrations

import (
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/migrate"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
)

func TestMain(m *testing.M) {
	def.Init()
	check.TestMain(m)
}

func Test(tt *testing.T) {
	t := check.T(tt)
	ctx := def.NewContext(app.ServiceName)
	migrate.UpDownTest(t, ctx, goose, ".", def.TestMySQLCfg(def.MySQLAuth{
		User: def.ExampleDBUser,
		Pass: def.ExampleDBPass,
		DB:   def.ExampleDBName,
	}))
}
