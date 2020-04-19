// Package api implements JSON-RPC 2.0 method handlers.
package api

import (
	"context"

	"github.com/powerman/go-monolith-example/internal/apiauth"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// API implements JSON-RPC 2.0 method handlers.
type API struct {
	a         app.Appl
	authn     apiauth.Authenticator
	strictErr bool
}

// New creates new net/rpc service.
func New(a app.Appl, authn apiauth.Authenticator) *API {
	return &API{
		a:     a,
		authn: authn,
	}
}
