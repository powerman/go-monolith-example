package jsonrpc2

import (
	"context"
	"errors"
	"io"
	"net/rpc"

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/rpc-codec/jsonrpc2"

	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/internal/apix"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/pkg/jsonrpc2x"
	"github.com/powerman/go-monolith-example/pkg/repo"
)

// APIErr converts non-JSON-RPC 2.0 errors to JSON-RPC 2.0 errors.
func apiErr(next jsonrpc2x.Handler) jsonrpc2x.Handler {
	return func() error {
		err := next()

		switch {
		case errors.Is(err, app.ErrAccessDenied):
			err = api.ErrForbidden
		case errors.Is(err, app.ErrNotFound):
			err = api.ErrNotFound
		// Common errors.
		case errors.Is(err, apix.ErrAccessTokenInvalid):
			err = api.ErrUnauthorized
		case errors.Is(err, context.DeadlineExceeded):
			err = api.ErrTryAgainLater
		case errors.Is(err, context.Canceled):
			err = api.ErrTryAgainLater
		case errors.Is(err, io.ErrUnexpectedEOF):
			err = api.ErrTryAgainLater
		case errors.Is(err, rpc.ErrShutdown):
			err = api.ErrTryAgainLater
		case errors.As(err, new(*mysql.MySQLError)):
			err = jsonrpc2x.ErrInternal
		case errors.Is(err, repo.ErrSchemaVer):
			err = jsonrpc2x.ErrInternal
		}

		if err == nil || errors.As(err, new(*jsonrpc2.Error)) {
			return err
		}
		return jsonrpc2.NewError(jsonrpc2x.ErrInternal.Code, err.Error())
	}
}
