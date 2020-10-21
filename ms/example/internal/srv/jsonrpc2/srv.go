// Package jsonrpc2 implements JSON-RPC 2.0 method handlers.
package jsonrpc2

import (
	"context"

	"github.com/powerman/go-monolith-example/internal/apix"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Config contains configuration for JSON-RPC 2.0 service.
type Config struct {
	StrictErr bool // Set to true to panic if RPC method returns undocumented error.
}

// Server implements JSON-RPC 2.0 method handlers.
type Server struct {
	appl  app.Appl
	authn apix.Authn
	cfg   Config
}

// New creates new net/rpc service.
func New(appl app.Appl, authn apix.Authn, cfg Config) *Server {
	return &Server{
		appl:  appl,
		authn: authn,
		cfg:   cfg,
	}
}
