package grpc

import (
	"errors"
	"fmt"

	"github.com/powerman/sensitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/powerman/go-monolith-example/api/proto/powerman/example/auth"
	"github.com/powerman/go-monolith-example/internal/apix"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
)

func (srv *server) CreateAccount(ctx Ctx, req *api.CreateAccountRequest) (*api.Account, error) {
	userID := req.GetAccountId()
	password := sensitive.String(req.GetAccount().GetPassword())
	user := app.User{
		Email:       req.GetAccount().GetEmail(),
		DisplayName: req.GetAccount().GetUser().GetDisplayName(),
	}

	err := srv.appl.Register(ctx, userID, password, &user)
	if err != nil {
		return nil, apiErr(err)
	}
	return apiAccount(user), nil
}

func (srv *server) SigninIdentity(ctx Ctx, req *api.SigninIdentityRequest) (_ *api.SigninIdentityResponse, err error) {
	var accessToken app.AccessToken
	switch req.Auth.(type) {
	default:
		panic(fmt.Sprintf("unknown req.Auth type: %T", req.Auth))
	case nil:
		return nil, status.Error(codes.InvalidArgument, "auth required")
	case *api.SigninIdentityRequest_Account:
		auth := req.GetAccount()
		accessToken, err = srv.appl.LoginByUserID(ctx, auth.GetAccountId(), sensitive.String(auth.GetPassword()))
	case *api.SigninIdentityRequest_Email:
		auth := req.GetEmail()
		accessToken, err = srv.appl.LoginByEmail(ctx, auth.GetEmail(), sensitive.String(auth.GetPassword()))
	}
	if err != nil {
		return nil, apiErr(err)
	}

	user, err := srv.appl.Authenticate(ctx, accessToken)
	switch {
	case errors.Is(err, app.ErrNotFound):
		return nil, status.Error(codes.Aborted, "access_token was deleted")
	case err != nil:
		return nil, apiErr(err)
	}

	resp := &api.SigninIdentityResponse{
		AccessToken: string(accessToken),
		User:        apiUser(*user),
	}
	return resp, nil
}

func (srv *server) SignoutIdentity(ctx Ctx, req *api.SignoutIdentityRequest) (_ *api.SignoutIdentityResponse, err error) {
	if req.GetEverywhere() {
		_, _, auth := apix.FromContext(ctx)
		err = srv.appl.LogoutUser(ctx, auth.UserName)
	} else {
		accessToken := app.AccessToken(apix.AccessTokenFromContext(ctx))
		err = srv.appl.Logout(ctx, accessToken)
	}
	if err != nil {
		return nil, apiErr(err)
	}
	return &api.SignoutIdentityResponse{}, nil
}

func (srv *server) CheckAccessToken(ctx Ctx, req *api.CheckAccessTokenRequest) (*api.CheckAccessTokenResponse, error) {
	accessToken := app.AccessToken(apix.AccessTokenFromContext(ctx))

	user, err := srv.appl.Authenticate(ctx, accessToken)
	switch {
	case errors.Is(err, app.ErrNotFound):
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	case err != nil:
		return nil, apiErr(err)
	}

	resp := &api.CheckAccessTokenResponse{
		User: apiUser(*user),
	}
	return resp, nil
}
