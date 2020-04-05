// Package migrations provides goose migrations for microservice.
package migrations

import (
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	goosepkg "github.com/powerman/goose/v2"
)

//nolint:gochecknoglobals // Force auto-generated to use instance.
var goose = def.NewGoose(app.ServiceName)

// Goose returns goose instance with loaded Go migrations from this dir.
func Goose() *goosepkg.Instance { return goose }
