// Package serve provides helpers to start and shutdown network services.
package serve

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
