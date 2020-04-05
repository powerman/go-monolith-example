package api

import (
	"errors"

	"github.com/powerman/go-monolith-example/internal/apiauth"
	"github.com/powerman/go-monolith-example/internal/jsonrpc2x"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/proto/rpc"
)

func makeAPIError(errsExtra []error) jsonrpc2x.Middleware {
	return func(next jsonrpc2x.Handler) jsonrpc2x.Handler {
		return func() error {
			err := next()

			switch {
			case err == nil:
			case errors.Is(err, rpc.ErrInvalidParams):
			case errors.Is(err, app.ErrAccessDenied):
				err = rpc.ErrForbidden
			case errors.Is(err, apiauth.ErrInvalidAccessToken):
				err = rpc.ErrUnauthorized
			case errors.Is(err, app.ErrNotFound):
				err = rpc.ErrNotFound
			}
			return sanitizeErr(err, errsExtra)
		}
	}
}

func sanitizeErr(err error, appErrs []error) error {
	if err == nil {
		return nil
	}
	for i := range rpc.ErrsCommon {
		if errors.Is(err, rpc.ErrsCommon[i]) {
			return err
		}
	}
	for i := range appErrs {
		if errors.Is(err, appErrs[i]) {
			return err
		}
	}
	return rpc.ErrInternal
}
