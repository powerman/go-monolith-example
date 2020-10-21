package jsonrpc2

import (
	"errors"

	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/internal/jsonrpc2x"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
)

func apiErr(next jsonrpc2x.Handler) jsonrpc2x.Handler { //nolint:gocyclo // By design.
	return func() error {
		err := next()

		switch {
		case errors.Is(err, app.ErrAccessDenied):
			err = api.ErrForbidden
		case errors.Is(err, app.ErrNotFound):
			err = api.ErrNotFound
		}

		return err
	}
}
