package grpc

import (
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/powerman/go-monolith-example/internal/apix"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
)

var (
	errRequireAuthn       = status.Error(codes.Unauthenticated, "authentication required")
	errInvalidAccessToken = status.Error(codes.Unauthenticated, "invalid access token")
)

func (srv *server) authn(ctx Ctx, fullMethod string) (Ctx, error) {
	ctx, auth, err := apix.GRPCNewContext(ctx, fullMethod, srv)

	switch {
	case strings.Contains(fullMethod, "/grpc."):
		return ctx, nil
	case strings.Contains(fullMethod, ".NoAuthSvc/"):
		return ctx, nil
	case errors.Is(err, apix.ErrAccessTokenInvalid):
		return ctx, errInvalidAccessToken
	case err != nil:
		return ctx, status.Error(codes.Internal, err.Error())
	case auth.UserName == dom.NoUser:
		return ctx, errRequireAuthn
	default:
		return ctx, nil
	}
}

// Authenticate implements apix.Authn.
func (srv *server) Authenticate(ctx Ctx, accessToken apix.AccessToken) (auth dom.Auth, err error) {
	user, err := srv.appl.Authenticate(ctx, app.AccessToken(accessToken))
	switch {
	case errors.Is(err, app.ErrNotFound):
		err = apix.ErrAccessTokenInvalid
	case err == nil:
		auth = dom.Auth{
			UserName: user.Name,
			Admin:    user.Role == app.RoleAdmin,
		}
	}
	return auth, err
}
