package def

import (
	"github.com/powerman/goose/v2"
	"github.com/powerman/structlog"
)

// NewGoose creates a goose instance with configured logger.
func NewGoose(service string) *goose.Instance {
	log := structlog.New(structlog.KeyApp, service, structlog.KeyUnit, "goose").
		SetKeysFormat(map[string]string{
			structlog.KeyMessage: " %[2]s",
		})
	g := goose.NewInstance()
	g.SetLogger(log)
	return g
}
