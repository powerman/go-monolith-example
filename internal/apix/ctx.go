package apix

import (
	"context"

	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/powerman/structlog"
	"github.com/sebest/xff"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/reflectx"
)

type contextKey int

const (
	_ contextKey = iota
	contextKeyRemote
	contextKeyMethodName
)

// Ctx describe Ctx param used by all API methods.
type Ctx struct {
	AccessToken AccessToken
	jsonrpc2.Ctx
}

// NewContext returns a new context.Context that carries values describing
// this request without any deadline, plus some of the values.
func (c *Ctx) NewContext(
	authn Authn,
	service string,
) (
	ctx context.Context,
	log *structlog.Logger,
	methodName string,
	auth dom.Auth,
	err error,
) {
	ctx = c.Context()

	remote := "unknown" // non-HTTP RPC call (like in tests)
	if r := jsonrpc2.HTTPRequestFromContext(ctx); r != nil {
		remote = xff.GetRemoteAddr(r)
	}
	ctx = context.WithValue(ctx, contextKeyRemote, remote)

	methodName = reflectx.CallerMethodName(1)
	ctx = context.WithValue(ctx, contextKeyMethodName, methodName)

	if c.AccessToken != "" {
		auth, err = authn.Authenticate(c.AccessToken)
	}

	log = structlog.New(
		structlog.KeyApp, service,
		def.LogRemote, remote,
		def.LogFunc, methodName,
		def.LogUserID, auth.UserID,
	)
	ctx = structlog.NewContext(ctx, log)

	c.SetContext(ctx)
	return ctx, log, methodName, auth, err
}

// FromContext returns values describing request stored in ctx, if any.
func FromContext(ctx context.Context) (remote, methodName string) {
	remote, _ = ctx.Value(contextKeyRemote).(string)
	methodName, _ = ctx.Value(contextKeyMethodName).(string)
	return remote, methodName
}
