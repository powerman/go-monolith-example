package apix

import (
	"context"
	"net"

	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/powerman/structlog"
	"github.com/sebest/xff"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/reflectx"
)

// JSONRPC2Ctx describe JSON-RPC 2.0 Ctx param used by all API methods.
type JSONRPC2Ctx struct {
	AccessToken string
	jsonrpc2.Ctx
}

// NewContext returns a new context.Context that carries values describing
// this request without any deadline, plus some of the values.
func (c *JSONRPC2Ctx) NewContext(
	authn Authn,
	service string,
) (
	ctx Ctx,
	log *structlog.Logger,
	methodName string,
	auth dom.Auth,
	err error,
) {
	ctx = c.Context()

	remoteIP := "" // non-HTTP RPC call (like in tests)
	if r := jsonrpc2.HTTPRequestFromContext(ctx); r != nil {
		remoteIP, _, _ = net.SplitHostPort(xff.GetRemoteAddr(r))
	}
	ctx = context.WithValue(ctx, contextKeyRemoteIP, remoteIP)

	methodName = reflectx.CallerMethodName(1)
	ctx = context.WithValue(ctx, contextKeyMethodName, methodName)

	if c.AccessToken != "" {
		accessToken := AccessToken(c.AccessToken)
		ctx = context.WithValue(ctx, contextKeyAccessToken, accessToken)
		auth, err = authn.Authenticate(ctx, accessToken)
		ctx = context.WithValue(ctx, contextKeyAuth, auth)
	}

	log = structlog.New(
		structlog.KeyApp, service,
		def.LogRemoteIP, remoteIP,
		def.LogFunc, methodName,
		def.LogUserName, auth.UserName,
	)
	ctx = structlog.NewContext(ctx, log)

	c.SetContext(ctx)
	return ctx, log, methodName, auth, err
}
