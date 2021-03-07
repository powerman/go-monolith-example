// +build integration

package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"testing"

	"github.com/parnurzeal/gorequest"
	"github.com/powerman/check"
	"github.com/powerman/pqx"
	"golang.org/x/oauth2"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/status"

	api "github.com/powerman/go-monolith-example/api/proto/powerman/example/auth"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

type (
	gwSigninIdentityResponse struct {
		AccessToken string
		User        gwUser
	}
	gwUser struct {
		Name        string
		DisplayName string
		Access      gwAccess
	}
	gwAccess struct {
		Role string
	}
)

func TestSmoke(tt *testing.T) {
	t := check.T(tt)

	s := &Service{cfg: cfg}

	pc, _, _, _ := runtime.Caller(0)
	suffix := runtime.FuncForPC(pc).Name()
	suffix = suffix[:strings.LastIndex(suffix, ".")]
	_, cleanup, err := pqx.EnsureTempDB(tLogger(*t), suffix, cfg.Postgres.Config)
	t.Must(t.Nil(err))
	defer cleanup()
	cfg.Postgres.DBName += "_" + suffix // Assign to cfg and not s.cfg as a reminder: they are the same.

	ctxStartup, cancel := context.WithTimeout(ctx, def.TestTimeout)
	defer cancel()
	ctxShutdown, shutdown := context.WithCancel(ctx)
	errc := make(chan error)
	go func() { errc <- s.RunServe(ctxStartup, ctxShutdown, shutdown) }()
	defer func() {
		shutdown()
		t.Nil(<-errc, "RunServe")
		if s.repo != nil {
			s.repo.Close()
		}
	}()
	t.Must(t.Nil(netx.WaitTCPPort(ctxStartup, cfg.BindAddr), "connect to gRPC service"))
	t.Must(t.Nil(netx.WaitTCPPort(ctxStartup, cfg.BindAddrInt), "connect to internal gRPC service"))
	t.Must(t.Nil(netx.WaitTCPPort(ctxStartup, cfg.BindGRPCGWAddr), "connect to grpc-gateway service"))

	ca, err := netx.LoadCACert(cfg.TLSCACert)
	t.Must(t.Nil(err))
	conn, err := grpcpkg.DialContext(ctx, cfg.BindAddr.String(),
		grpcpkg.WithTransportCredentials(credentials.NewClientTLSFromCert(ca, "")),
		grpcpkg.WithBlock(),
	)
	t.Must(t.Nil(err, "grpc.Dial"))
	connInt, err := grpcpkg.DialContext(ctx, cfg.BindAddrInt.String(),
		grpcpkg.WithTransportCredentials(credentials.NewClientTLSFromCert(ca, "")),
		grpcpkg.WithBlock(),
	)
	t.Must(t.Nil(err, "grpc.Dial internal"))
	clientNoAuth := api.NewNoAuthSvcClient(conn)
	clientAuth := api.NewAuthSvcClient(conn)
	clientAuthInt := api.NewAuthIntSvcClient(connInt)
	clientGW := gorequest.New().
		TLSClientConfig(&tls.Config{RootCAs: ca}).
		Timeout(def.TestTimeout).
		Retry(30, def.TestSecond/10, http.StatusServiceUnavailable)
	endpointGW := fmt.Sprintf("https://%s/", s.cfg.BindGRPCGWAddr)

	var (
		userAdmin = &api.User{
			Name: "users/admin",
			Access: &api.Access{
				Role: api.Access_ROLE_ADMIN,
			},
		}
		userUser = &api.User{
			DisplayName: "U.S.E.R.",
			Access: &api.Access{
				Role: api.Access_ROLE_USER,
			},
		}
		authAdmin        grpcpkg.CallOption
		authUser         grpcpkg.CallOption
		accessTokenAdmin string
	)

	{ // register admin
		res, err := clientNoAuth.CreateAccount(ctx, &api.CreateAccountRequest{
			AccountId: "admin",
		})
		t.Nil(err)
		t.DeepEqual(res, &api.Account{
			Name:       "accounts/admin",
			User:       userAdmin,
			CreateTime: res.CreateTime,
		})
	}
	{ // register user
		res, err := clientNoAuth.CreateAccount(ctx, &api.CreateAccountRequest{
			Account: &api.Account{
				User: &api.User{
					DisplayName: "U.S.E.R.",
					Access: &api.Access{
						Role: api.Access_ROLE_ADMIN,
					},
				},
				Password: "pass",
				Email:    "user@example.com",
			},
		})
		t.Nil(err)
		userIDUser := strings.TrimPrefix(res.Name, "accounts/")
		userUser.Name = "users/" + userIDUser
		t.DeepEqual(res, &api.Account{
			Name:       "accounts/" + userIDUser,
			User:       userUser,
			Email:      "user@example.com",
			CreateTime: res.CreateTime,
		})
	}
	{ // login admin
		res, err := clientNoAuth.SigninIdentity(ctx, &api.SigninIdentityRequest{
			Auth: &api.SigninIdentityRequest_Account{Account: &api.SigninIdentityRequest_AccountAuth{
				AccountId: "admin",
			}},
		})
		t.Nil(err)
		authAdmin = grpcpkg.PerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{
			AccessToken: res.AccessToken,
		}))
		t.DeepEqual(res.User, userAdmin)
	}
	{ // login user
		res, err := clientNoAuth.SigninIdentity(ctx, &api.SigninIdentityRequest{
			Auth: &api.SigninIdentityRequest_Email{Email: &api.SigninIdentityRequest_EmailAuth{
				Email:    "user@example.com",
				Password: "pass",
			}},
		})
		t.Nil(err)
		authUser = grpcpkg.PerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{
			AccessToken: res.AccessToken,
		}))
		t.DeepEqual(res.User, userUser)
	}
	{ // authenticate
		res, err := clientAuthInt.CheckAccessToken(ctx, &api.CheckAccessTokenRequest{}, authAdmin)
		t.Nil(err)
		t.DeepEqual(res.User, userAdmin)
	}
	{ // logout
		_, err := clientAuth.SignoutIdentity(ctx, &api.SignoutIdentityRequest{}, authUser)
		t.Nil(err)
	}
	{ // authenticate
		res, err := clientAuthInt.CheckAccessToken(ctx, &api.CheckAccessTokenRequest{}, authUser)
		t.Err(err, status.Error(codes.Unauthenticated, "invalid access token"))
		t.Nil(res)
	}
	{ // grpc-gateway login
		var body gwSigninIdentityResponse
		resp, _, errs := clientGW.Clone().Post(endpointGW + "accounts:signinIdentity").
			Send(`{"account":{"accountId":"admin"}}`).
			EndStruct(&body)
		t.Nil(errs)
		t.Equal(resp.StatusCode, 200)
		t.Len(body.AccessToken, 26)
		t.DeepEqual(body.User, gwUser{
			Name: "users/admin",
			Access: gwAccess{
				Role: "ROLE_ADMIN",
			},
		})
		accessTokenAdmin = body.AccessToken
	}
	{ // grpc-gateway logout
		resp, _, errs := clientGW.Clone().Post(endpointGW+"accounts:signoutIdentity").
			Set("Authorization", "bearer "+accessTokenAdmin).
			Send(`{}`).
			End()
		t.Nil(errs)
		t.Equal(resp.StatusCode, 200)
		resp, _, errs = clientGW.Clone().Post(endpointGW+"accounts:signoutIdentity").
			Set("Authorization", "bearer "+accessTokenAdmin).
			Send(`{}`).
			End()
		t.Nil(errs)
		t.Equal(resp.StatusCode, 401)
	}
}
