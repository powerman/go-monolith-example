package rpc

import (
	"context"

	"github.com/powerman/go-monolith-example/internal/apiauth"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/internal/reflectx"
	"github.com/powerman/go-monolith-example/proto"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/powerman/structlog"
	"github.com/sebest/xff"
)

type contextKey int

const (
	_ contextKey = iota
	contextKeyRemote
	contextKeyMethodName
)

// Ctx describe Ctx param used by all API methods.
type Ctx struct {
	AccessToken proto.AccessToken
	jsonrpc2.Ctx
}

// NewContext returns a new Context that carries values describing this
// RPC request without any deadline, plus some of these values.
func (c *Ctx) NewContext(service string) (context.Context, *structlog.Logger, string, dom.Auth, error) {
	ctx := c.Context()
	if ctx == nil {
		ctx = context.Background() // non-HTTP RPC call (like in tests)
	}

	remote := "unknown" // non-HTTP RPC call (like in tests)
	if r := jsonrpc2.HTTPRequestFromContext(ctx); r != nil {
		remote = xff.GetRemoteAddr(r)
	}
	ctx = context.WithValue(ctx, contextKeyRemote, remote)

	methodName := reflectx.CallerMethodName(1)
	ctx = context.WithValue(ctx, contextKeyMethodName, methodName)

	auth, err := apiauth.ParseAccessToken(c.AccessToken)

	log := structlog.New(
		structlog.KeyApp, service,
		def.LogRemote, remote,
		def.LogFunc, methodName,
		def.LogUser, auth.UserID,
	)
	ctx = structlog.NewContext(ctx, log)

	return ctx, log, methodName, auth, err
}

// FromContext returns values describing RPC request stored in ctx, if any.
func FromContext(ctx context.Context) (remote, methodName string) {
	remote, _ = ctx.Value(contextKeyRemote).(string)
	methodName, _ = ctx.Value(contextKeyMethodName).(string)
	return remote, methodName
}
