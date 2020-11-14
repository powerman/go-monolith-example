// +build integration

package migrations_test

import (
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/ms/example/internal/config"
	"github.com/powerman/go-monolith-example/ms/example/internal/migrations"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/migrate"
)

var cfg *config.ServeConfig

func TestMain(m *testing.M) {
	def.Init()
	cfg = config.MustGetServeTest()
	check.TestMain(m)
}

func Test(tt *testing.T) {
	t := check.T(tt)
	ctx := def.NewContext(app.ServiceName)
	migrate.MySQLUpDownTest(t, ctx, migrations.Goose(), ".", cfg.MySQL)
}
