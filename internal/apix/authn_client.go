package apix

import (
	"context"
	"crypto/x509"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/powerman/go-monolith-example/api/proto/powerman/example/auth"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/pkg/grpcx"
)

type authnClient struct {
	client api.AuthIntSvcClient
}

// NewAuthnClient returns Authn implementation using gRPC AuthIntSvc at addr.
func NewAuthnClient(
	ctx Ctx,
	reg *prometheus.Registry,
	service string,
	ca *x509.CertPool,
	addr string,
) (Authn, error) {
	const subsystem = "apix"
	metrics := grpcx.NewClientMetrics(reg, service, subsystem)
	conn, err := grpcx.Dial(ctx, addr, service, metrics, ca)
	if err != nil {
		return nil, err
	}
	client := api.NewAuthIntSvcClient(conn)
	return &authnClient{client: client}, nil
}

func (c *authnClient) Authenticate(ctx Ctx, accessToken AccessToken) (auth dom.Auth, err error) {
	const rpcTimeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, rpcTimeout)
	defer cancel()

	creds := grpcx.AccessTokenCreds(string(accessToken))
	resp, err := c.client.CheckAccessToken(ctx, &api.CheckAccessTokenRequest{}, creds)
	var userName *dom.UserName
	if err == nil {
		userName, err = dom.ParseUserName(resp.GetUser().GetName())
	}
	if err == nil {
		auth.UserName = *userName
		auth.Admin = resp.GetUser().GetAccess().GetRole() == api.Access_ROLE_ADMIN
	}
	if status.Code(err) == codes.Unauthenticated {
		err = ErrAccessTokenInvalid
	}
	return auth, err
}
