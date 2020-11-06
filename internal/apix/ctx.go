package apix

import (
	"context"

	"github.com/powerman/go-monolith-example/internal/dom"
)

type Ctx = context.Context

type contextKey int

const (
	_ contextKey = iota
	contextKeyRemote
	contextKeyMethodName
	contextKeyAuth
	contextKeyAccessToken
)

// FromContext returns values describing request stored in ctx, if any.
func FromContext(ctx Ctx) (remote, methodName string, auth dom.Auth) {
	remote, _ = ctx.Value(contextKeyRemote).(string)
	methodName, _ = ctx.Value(contextKeyMethodName).(string)
	auth, _ = ctx.Value(contextKeyAuth).(dom.Auth)
	return remote, methodName, auth
}

// AccessTokenFromContext returns AccessToken stored in ctx, if any.
func AccessTokenFromContext(ctx Ctx) (accessToken AccessToken) {
	accessToken, _ = ctx.Value(contextKeyAccessToken).(AccessToken)
	return
}
