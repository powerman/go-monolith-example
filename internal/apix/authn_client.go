package apix

import (
	"crypto/x509"

	"github.com/prometheus/client_golang/prometheus"

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
	resp, err := c.client.CheckAccessToken(ctx, nil, grpcx.Token(string(accessToken)))
	var userName *dom.UserName
	if err == nil {
		userName, err = dom.ParseUserName(resp.GetUser().GetName())
	}
	if err == nil {
		auth.UserName = *userName
		auth.Admin = resp.GetUser().GetAccess().GetRole() == api.Access_ROLE_ADMIN
	}
	return auth, err
}
