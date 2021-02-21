// Package grpcgw provides grpc-gateway server.
package grpcgw

import (
	"context"
	"crypto/x509"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/laher/mergefs"

	"github.com/powerman/go-monolith-example/api/proto"
	api "github.com/powerman/go-monolith-example/api/proto/powerman/example/auth"
	"github.com/powerman/go-monolith-example/pkg/grpcx"
	"github.com/powerman/go-monolith-example/pkg/netx"
	"github.com/powerman/go-monolith-example/third_party"
	"github.com/powerman/go-monolith-example/web"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Config contains configuration for grpc-gateway server.
type Config struct {
	CtxShutdown      Ctx
	Endpoint         netx.Addr
	CA               *x509.CertPool
	GRPCGWPattern    string // Pattern for http.ServeMux to serve grpc-gateway.
	OpenAPIPattern   string // Pattern for http.ServeMux to serve swagger.json.
	SwaggerUIPattern string // Pattern for http.ServeMux to serve Swagger UI.
}

// NewServer creates and returns HTTP handler with grpc-gateway service on
// given cfg.Pattern.
func NewServer(cfg Config) (*http.ServeMux, error) {
	ctx, addr, opts := cfg.CtxShutdown, cfg.Endpoint.String(), grpcx.DialOptions(cfg.CA)

	gw := runtime.NewServeMux()
	err := api.RegisterNoAuthSvcHandlerFromEndpoint(ctx, gw, addr, opts)
	if err == nil {
		err = api.RegisterAuthSvcHandlerFromEndpoint(ctx, gw, addr, opts)
	}
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle(cfg.GRPCGWPattern, noCache(corsAllowAll(gw)))
	mux.Handle(cfg.OpenAPIPattern, http.StripPrefix(cfg.OpenAPIPattern, http.FileServer(http.FS(proto.OpenAPI))))
	mux.Handle(cfg.SwaggerUIPattern, http.StripPrefix(cfg.SwaggerUIPattern, http.FileServer(http.FS(
		mergefs.Merge(web.SwaggerUI, third_party.SwaggerUI)))))
	return mux, nil
}
