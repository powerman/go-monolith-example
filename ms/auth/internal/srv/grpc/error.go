package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
)

// APIErr converts non-gRPC errors to gRPC errors.
func apiErr(err error) error {
	if err == nil {
		return nil
	}

	code := codes.Internal
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		code = codes.DeadlineExceeded
	case errors.Is(err, context.Canceled):
		code = codes.Canceled
	case errors.Is(err, app.ErrAccessDenied):
		code = codes.PermissionDenied
	case errors.Is(err, app.ErrNotFound):
		code = codes.NotFound
	case errors.Is(err, app.ErrAlreadyExist):
		code = codes.AlreadyExists
	case errors.Is(err, app.ErrWrongPassword):
		code = codes.NotFound
		err = app.ErrNotFound // Hide actual error.
	case errors.Is(err, app.ErrValidate):
		code = codes.InvalidArgument
	}
	return status.Error(code, err.Error())
}
