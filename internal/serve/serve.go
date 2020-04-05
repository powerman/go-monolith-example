// Package serve provides helpers to start and shutdown any services.
package serve

import (
	"context"
	"net/http"
	"net/rpc"

	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Addr is a serve addr (using tcp network).
type Addr interface {
	Host() string
	Port() int
	String() string
}

// HTTP starts HTTP server on addr using handler logged as service.
// It runs until failed or ctx.Done.
func HTTP(ctx Ctx, addr Addr, handler http.Handler, service string) error {
	log := structlog.FromContext(ctx, nil).New(def.LogService, service)

	srv := &http.Server{
		Addr:    addr.String(),
		Handler: handler,
	}

	log.Info("serve", def.LogHost, addr.Host(), def.LogPort, addr.Port())
	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServe() }()

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
func Metrics(ctx Ctx, addr Addr, reg *prometheus.Registry) error {
	handler := promhttp.InstrumentMetricHandler(reg, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	mux := http.NewServeMux()
	mux.Handle("/metrics", handler)
	return HTTP(ctx, addr, mux, "Prometheus metrics")
}

// RPC starts HTTP server on addr path /rpc using rcvr as JSON-RPC 2.0
// handler.
func RPC(ctx Ctx, addr Addr, rcvr interface{}) error {
	srv := rpc.NewServer()
	err := srv.Register(rcvr)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/rpc", jsonrpc2.HTTPHandler(srv))
	return HTTP(ctx, addr, mux, "JSON-RPC 2.0")
}
