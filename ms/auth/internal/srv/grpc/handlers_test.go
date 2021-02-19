package grpc_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/sensitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/powerman/go-monolith-example/api/proto/powerman/example/auth"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/pkg/grpcx"
)

var (
	now       = time.Now()
	userAdmin = app.User{
		Name:        dom.NewUserName("admin"),
		Email:       "root@host",
		DisplayName: "Root",
		Role:        app.RoleAdmin,
		CreateTime:  now,
	}
	apiUserAdmin = &api.User{
		Name:        userAdmin.Name.String(),
		DisplayName: userAdmin.DisplayName,
		Access: &api.Access{
			Role: api.Access_ROLE_ADMIN,
		},
	}
)

func TestCreateAccount(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	clientNoAuth, _, _, mockAppl := testNew(t)

	var (
		userRand1 = app.User{
			Name:       dom.NewUserName("random-id-1"),
			Email:      "",
			Role:       app.RoleUser,
			CreateTime: now,
		}
		accRand1 = &api.Account{
			Name: "accounts/random-id-1",
			User: &api.User{
				Name:        "users/random-id-1",
				DisplayName: "",
				Access: &api.Access{
					Role: api.Access_ROLE_USER,
				},
			},
			Email:      "",
			CreateTime: timestamppb.New(now),
		}
		accAdmin = &api.Account{
			Name: "accounts/admin",
			User: &api.User{
				Name:        "users/admin",
				DisplayName: "Root",
				Access: &api.Access{
					Role: api.Access_ROLE_ADMIN,
				},
			},
			Email:      "root@host",
			CreateTime: timestamppb.New(now),
		}
	)

	mockAppl.EXPECT().Register(gomock.Any(), "", sensitive.String(""),
		app.MatchUser{User: &app.User{}}).
		SetArg(3, userRand1)
	mockAppl.EXPECT().Register(gomock.Any(), "admin", sensitive.String("secret"),
		app.MatchUser{User: &app.User{
			Email:       "root@host",
			DisplayName: "Root",
		}}).
		SetArg(3, userAdmin)
	mockAppl.EXPECT().Register(gomock.Any(), "admin", sensitive.String(""), app.MatchUser{User: &app.User{}}).
		Return(app.ErrAlreadyExist)
	mockAppl.EXPECT().Register(gomock.Any(), "bad", sensitive.String(""), app.MatchUser{User: &app.User{}}).
		Return(fmt.Errorf("%w: userID", app.ErrValidate))

	tests := []struct {
		accountID   string
		password    string
		email       string
		displayName string
		want        *api.Account
		wantCode    codes.Code
		wantErr     string
	}{
		{"", "", "", "", accRand1, codes.OK, ``},
		{"admin", "secret", "root@host", "Root", accAdmin, codes.OK, ``},
		{"admin", "", "", "", nil, codes.AlreadyExists, `already exists`},
		{"bad", "", "", "", nil, codes.InvalidArgument, `validate: userID`},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			arg := &api.CreateAccountRequest{
				Account: &api.Account{
					User: &api.User{
						DisplayName: tc.displayName,
					},
					Password: tc.password,
					Email:    tc.email,
				},
				AccountId: tc.accountID,
			}
			msg, err := clientNoAuth.CreateAccount(ctx, arg)
			if tc.wantCode == codes.OK {
				t.Nil(err)
				t.DeepEqual(msg, tc.want)
			} else {
				t.Err(err, status.Error(tc.wantCode, tc.wantErr))
				t.Nil(msg)
			}
		})
	}
}

func TestSigninIdentity(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	clientNoAuth, _, _, mockAppl := testNew(t)

	var (
		respAdmin1 = &api.SigninIdentityResponse{
			AccessToken: "token1",
			User:        apiUserAdmin,
		}
		respAdmin2 = &api.SigninIdentityResponse{
			AccessToken: "token2",
			User:        apiUserAdmin,
		}
	)

	mockAppl.EXPECT().LoginByUserID(gomock.Any(), "admin", sensitive.String("secret")).Return(app.AccessToken("token1"), nil)
	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("token1")).Return(&userAdmin, nil)
	mockAppl.EXPECT().LoginByEmail(gomock.Any(), "root@host", sensitive.String("secret")).Return(app.AccessToken("token2"), nil)
	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("token2")).Return(&userAdmin, nil)
	mockAppl.EXPECT().LoginByEmail(gomock.Any(), "root@host", sensitive.String("secret")).Return(app.AccessToken("token3"), nil)
	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("token3")).Return(nil, app.ErrNotFound)
	mockAppl.EXPECT().LoginByUserID(gomock.Any(), "admin", gomock.Any()).Return(app.AccessToken(""), app.ErrWrongPassword)
	mockAppl.EXPECT().LoginByEmail(gomock.Any(), "user@host", sensitive.String("")).Return(app.AccessToken(""), app.ErrNotFound)

	tests := []struct {
		auth     interface{}
		want     *api.SigninIdentityResponse
		wantCode codes.Code
		wantErr  string
	}{
		{nil, nil, codes.InvalidArgument, "auth required"},
		{&api.SigninIdentityRequest_AccountAuth{
			AccountId: "admin",
			Password:  "secret",
		}, respAdmin1, codes.OK, ""},
		{&api.SigninIdentityRequest_EmailAuth{
			Email:    "root@host",
			Password: "secret",
		}, respAdmin2, codes.OK, ""},
		{&api.SigninIdentityRequest_EmailAuth{
			Email:    "root@host",
			Password: "secret",
		}, nil, codes.Aborted, "access_token was deleted"},
		{&api.SigninIdentityRequest_AccountAuth{
			AccountId: "admin",
			Password:  "",
		}, nil, codes.NotFound, "not found"},
		{&api.SigninIdentityRequest_EmailAuth{
			Email:    "user@host",
			Password: "",
		}, nil, codes.NotFound, "not found"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			req := &api.SigninIdentityRequest{}
			switch auth := tc.auth.(type) {
			case *api.SigninIdentityRequest_AccountAuth:
				req.Auth = &api.SigninIdentityRequest_Account{Account: auth}
			case *api.SigninIdentityRequest_EmailAuth:
				req.Auth = &api.SigninIdentityRequest_Email{Email: auth}
			}
			msg, err := clientNoAuth.SigninIdentity(ctx, req)
			if tc.wantCode == codes.OK {
				t.Nil(err)
				t.DeepEqual(msg, tc.want)
			} else {
				t.Err(err, status.Error(tc.wantCode, tc.wantErr))
				t.Nil(msg)
			}
		})
	}
}

func TestSignoutIdentity(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	_, clientAuth, _, mockAppl := testNew(t)

	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("boom")).Return(nil, io.EOF)
	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("token1")).Return(&userAdmin, nil)
	mockAppl.EXPECT().Logout(gomock.Any(), app.AccessToken("token1")).Return(nil)
	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("token2")).Return(&userAdmin, nil)
	mockAppl.EXPECT().LogoutUser(gomock.Any(), userAdmin.Name).Return(nil)
	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("token3")).Return(&userAdmin, nil)
	mockAppl.EXPECT().LogoutUser(gomock.Any(), userAdmin.Name).Return(io.EOF)
	mockAppl.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound).AnyTimes()

	tests := []struct {
		everywhere bool
		token      string
		wantCode   codes.Code
		wantErr    string
	}{
		// Test auth interceptor:
		{false, "", codes.Unauthenticated, "authentication required"},
		{false, "wrong", codes.Unauthenticated, "invalid access token"},
		{false, "boom", codes.Internal, "internal error"},
		// Test SignoutIdentity:
		{false, "token1", codes.OK, ""},
		{true, "token2", codes.OK, ""},
		{true, "token3", codes.Internal, "internal error"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			arg := &api.SignoutIdentityRequest{Everywhere: tc.everywhere}
			_, err := clientAuth.SignoutIdentity(ctx, arg, grpcx.AccessTokenCreds(tc.token))
			if tc.wantCode == codes.OK {
				t.Nil(err)
			} else {
				t.Err(err, status.Error(tc.wantCode, tc.wantErr))
			}
		})
	}
}

func TestCheckAccessToken(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	_, _, clientAuthInt, mockAppl := testNew(t)

	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("expire")).Return(&userAdmin, nil)
	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("expire")).Return(nil, app.ErrNotFound)
	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("boom")).Return(&userAdmin, nil)
	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("boom")).Return(nil, io.EOF)
	mockAppl.EXPECT().Authenticate(gomock.Any(), app.AccessToken("token1")).Return(&userAdmin, nil).Times(2)
	mockAppl.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound).AnyTimes()

	tests := []struct {
		token    string
		want     *api.CheckAccessTokenResponse
		wantCode codes.Code
		wantErr  string
	}{
		{"wrong", nil, codes.Unauthenticated, "invalid access token"},
		{"expire", nil, codes.Unauthenticated, "invalid access token"},
		{"boom", nil, codes.Internal, "internal error"},
		{"token1", &api.CheckAccessTokenResponse{User: apiUserAdmin}, codes.OK, ""},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			arg := &api.CheckAccessTokenRequest{}
			msg, err := clientAuthInt.CheckAccessToken(ctx, arg, grpcx.AccessTokenCreds(tc.token))
			if tc.wantCode == codes.OK {
				t.Nil(err)
				t.DeepEqual(msg, tc.want)
			} else {
				t.Err(err, status.Error(tc.wantCode, tc.wantErr))
				t.Nil(msg)
			}
		})
	}
}
