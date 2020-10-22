package def

import (
	"context"

	"github.com/powerman/structlog"
)

// NewContext returns context.Background() which contains logger
// configured for given service.
func NewContext(service string) context.Context {
	return structlog.NewContext(context.Background(), structlog.New(structlog.KeyApp, service))
}
