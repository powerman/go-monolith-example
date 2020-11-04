// Package jsonrpc2x provide helpers for JSON-RPC 2.0 API.
package jsonrpc2x

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/reflectx"
)

// Log is a synonym for convenience.
type Log = *structlog.Logger

// Handler is a JSON-RPC 2.0 method handler.
type Handler func() error

// Middleware is a JSON-RPC 2.0 middleware.
type Middleware func(Handler) Handler

// MakeValidateErr creates middleware which validates error against
// documented errors (api.ErrsCommon + api.ErrsExtra[method]).
//
// Use NewError instead of jsonrpc2.NewError to create errors which must
// match documented errors only by code.
//
// TODO Add new metric to report and extra (metric, methodName) args.
func MakeValidateErr(log Log, strict bool, errsCommon, errsExtra []error) Middleware { //nolint:gocognit // Questionable.
	log = log.New(structlog.KeyUnit, reflectx.CallerPkg(1))
	report := func(err error) {
		if strict {
			panic(err)
		} else {
			log.Warn(err)
		}
	}
	return func(next Handler) Handler {
		return func() error {
			err := next()

			if err == nil {
				return nil
			}
			for i := range errsCommon {
				if errors.Is(err, errsCommon[i]) {
					return err
				}
			}
			for i := range errsExtra {
				if errors.Is(err, errsExtra[i]) {
					return err
				}
			}
			if errors.As(err, new(*jsonrpc2.Error)) {
				report(fmt.Errorf("not documented (add to api.ErrsExtra): %w", err))
			} else {
				report(fmt.Errorf("not jsonrpc2.Error (add to srv/jsonrpc2.apiErr): %w", err))
			}
			return err
		}
	}
}

// MakeRecovery creates middleware which handle panics.
func MakeRecovery(log Log, metric def.Metrics) Middleware {
	log = log.New(structlog.KeyUnit, reflectx.CallerPkg(1))
	return func(next Handler) Handler {
		return func() (err error) {
			panicked := true
			defer func() {
				if p := recover(); panicked {
					err = ErrInternal
					metric.PanicsTotal.Inc()
					log.PrintErr("panic", "err", p, structlog.KeyStack, structlog.Auto)
				}
			}()
			err = next()
			panicked = false
			return err
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
			l = prometheus.Labels{
				methodLabel: methodName,
				failedLabel: strconv.FormatBool(err != nil),
			}
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
