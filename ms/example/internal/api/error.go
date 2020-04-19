package api

import (
	"errors"

	"github.com/powerman/go-monolith-example/internal/jsonrpc2x"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/proto/rpc"
)

func protoErr(next jsonrpc2x.Handler) jsonrpc2x.Handler {
	return func() error {
		err := next()

		switch {
		case errors.Is(err, app.ErrAccessDenied):
			err = rpc.ErrForbidden
		case errors.Is(err, app.ErrNotFound):
			err = rpc.ErrNotFound
		}

		return err
	}
}
