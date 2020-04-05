// Package jsonrpc2x provide helpers for JSON-RPC 2.0 API.
package jsonrpc2x

import (
	"time"

	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/reflectx"
	"github.com/powerman/go-monolith-example/proto/rpc"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
)

// Log is a synonym for convenience.
type Log = *structlog.Logger

// Handler is a JSON-RPC 2.0 method handler.
type Handler func() error

// Middleware is a JSON-RPC 2.0 middleware.
type Middleware func(Handler) Handler

// MakeRecovery creates middleware which handle panics.
func MakeRecovery(log Log, metric def.Metrics) Middleware {
	log = log.New(structlog.KeyUnit, reflectx.CallerPkg(1))
	return func(next Handler) Handler {
		return func() (err error) {
			defer func() {
				if p := recover(); p != nil {
					err = rpc.ErrInternal
					metric.PanicsTotal.Inc()
					log.PrintErr("panic", "err", p, structlog.KeyStack, structlog.Auto)
				}
			}()

			return next()
		}
	}
}

// MakeMetrics creates middleware which add default metrics.
func MakeMetrics(metric Metrics, methodName string) Middleware {
	return func(next Handler) Handler {
		return func() error {
			metric.reqInFlight.Inc()
			defer metric.reqInFlight.Dec()
			start := time.Now()

			err := next()

			l := prometheus.Labels{
				methodLabel: methodName,
				codeLabel:   code(err),
			}
			metric.reqTotal.With(l).Inc()
			metric.reqDuration.With(l).Observe(time.Since(start).Seconds())
			return err
		}
	}
}

// MakeAccessLog creates middleware which log method call success/failure.
func MakeAccessLog(log Log) Middleware {
	log = log.New(structlog.KeyUnit, reflectx.CallerPkg(1))
	return func(next Handler) Handler {
		return func() error {
			err := next()

			if err == nil {
				log.Info("handled")
			} else {
				log.PrintErr("failed to handle", "err", dropcode(err))
			}
			return err
		}
	}
}
