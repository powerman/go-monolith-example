// Package jsonrpc2 implements JSON-RPC 2.0 method handlers.
package jsonrpc2

import (
	"context"
	"net/http"
	"net/rpc"

	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/rs/cors"

	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/internal/apix"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Config contains configuration for JSON-RPC 2.0 service.
type Config struct {
	Pattern   string // Pattern for http.ServeMux.
	StrictErr bool   // Set to true to panic if RPC method returns undocumented error.
}

// Server implements JSON-RPC 2.0 method handlers.
type Server struct {
	appl  app.Appl
	authn apix.Authn
	cfg   Config
}

// NewServer creates and returns HTTP handler with JSON-RPC 2.0 service on
// given cfg.Pattern.
func NewServer(appl app.Appl, authn apix.Authn, cfg Config) *http.ServeMux {
	srv := &Server{
		appl:  appl,
		authn: authn,
		cfg:   cfg,
	}

	rpcSrv := rpc.NewServer()
	err := rpcSrv.RegisterName(api.Name, srv)
	if err != nil {
		panic(err)
	}

	handler := jsonrpc2.HTTPHandler(rpcSrv)
	handler = cors.AllowAll().Handler(handler)

	mux := http.NewServeMux()
	mux.Handle(srv.cfg.Pattern, handler)
	return mux
}
