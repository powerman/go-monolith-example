// Package serve provides helpers to start and shutdown network services.
package serve

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/rpc"

	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// HTTP starts HTTP server on addr using handler logged as service.
// It runs until failed or ctx.Done.
func HTTP(ctx Ctx, addr netx.Addr, tlsConfig *tls.Config, handler http.Handler, service string) error {
	log := structlog.FromContext(ctx, nil).New(def.LogServer, service)

	srv := &http.Server{
		Addr:      addr.String(),
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	log.Info("serve", def.LogHost, addr.Host(), def.LogPort, addr.Port())
	errc := make(chan error, 1)
	if srv.TLSConfig == nil {
		go func() { errc <- srv.ListenAndServe() }()
	} else {
		go func() { errc <- srv.ListenAndServeTLS("", "") }()
	}

	var err error
	select {
	case err = <-errc:
	case <-ctx.Done():
		err = srv.Shutdown(context.Background())
	}
	if err != nil {
		return log.Err("failed to serve", "err", err)
	}
	log.Info("shutdown")
	return nil
}

// Metrics starts HTTP server on addr path /metrics using reg as
// prometheus handler.
func Metrics(ctx Ctx, addr netx.Addr, reg *prometheus.Registry) error {
	mux := http.NewServeMux()
	HandleMetrics(mux, reg)
	return HTTP(ctx, addr, nil, mux, "Prometheus metrics")
}

// HandleMetrics adds reg's prometheus handler on /metrics at mux.
func HandleMetrics(mux *http.ServeMux, reg *prometheus.Registry) {
	handler := promhttp.InstrumentMetricHandler(reg, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	mux.Handle("/metrics", handler)
}

// RPC starts HTTP server on addr path /rpc using rcvr as JSON-RPC 2.0
// handler.
func RPC(ctx Ctx, addr netx.Addr, tlsConfig *tls.Config, rcvr interface{}) error {
	return RPCName(ctx, addr, tlsConfig, rcvr, "")
}

// RPCName starts HTTP server on addr path /rpc using rcvr as JSON-RPC 2.0
// handler but uses the provided name for the type instead of the
// receiver's concrete type.
func RPCName(ctx Ctx, addr netx.Addr, tlsConfig *tls.Config, rcvr interface{}, name string) (err error) {
	srv := rpc.NewServer()
	if name != "" {
		err = srv.RegisterName(name, rcvr)
	} else {
		err = srv.Register(rcvr)
	}
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/rpc", jsonrpc2.HTTPHandler(srv))
	return HTTP(ctx, addr, tlsConfig, mux, "JSON-RPC 2.0")
}

// GRPC starts gRPC server on addr, logged as service.
// It runs until failed or ctx.Done.
func GRPC(ctx Ctx, addr netx.Addr, srv *grpc.Server, service string) (err error) {
	log := structlog.FromContext(ctx, nil).New(def.LogServer, service)

	ln, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}

	log.Info("serve", def.LogHost, addr.Host(), def.LogPort, addr.Port())
	errc := make(chan error, 1)
	go func() { errc <- srv.Serve(ln) }()

	select {
	case err = <-errc:
	case <-ctx.Done():
		srv.GracefulStop() // It will not interrupt streaming.
	}
	if err != nil {
		return log.Err("failed to serve", "err", err)
	}
	log.Info("shutdown")
	return nil
}
