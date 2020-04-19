// Package jsonrpc2x provide helpers for JSON-RPC 2.0 API.
package jsonrpc2x

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/go-monolith-example/internal/apiauth"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/reflectx"
	"github.com/powerman/go-monolith-example/internal/repo"
	"github.com/powerman/go-monolith-example/proto/rpc"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
)

// Log is a synonym for convenience.
type Log = *structlog.Logger

// Handler is a JSON-RPC 2.0 method handler.
type Handler func() error

// Middleware is a JSON-RPC 2.0 middleware.
type Middleware func(Handler) Handler

// MakeValidateErr creates middleware which validates error against
// documented errors (rpc.ErrsCommon + proto.ErrsExtra[method]).
func MakeValidateErr(log Log, strict bool, errsExtra []error) Middleware { //nolint:gocognit // Questionable.
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
			for i := range rpc.ErrsCommon {
				if errors.Is(err, rpc.ErrsCommon[i]) {
					return err
				}
			}
			for i := range errsExtra {
				if errors.Is(err, errsExtra[i]) {
					return err
				}
			}
			if errors.As(err, new(*jsonrpc2.Error)) {
				report(fmt.Errorf("not documented (add to proto.ErrsExtra): %w", err))
			} else {
				report(fmt.Errorf("not jsonrpc2.Error (add to api.protoError): %w", err))
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

// ProtoErr converts common errors to JSON-RPC 2.0 errors.
func ProtoErr(next Handler) Handler {
	return func() error {
		err := next()

		switch {
		case errors.Is(err, apiauth.ErrAccessTokenInvalid):
			err = rpc.ErrUnauthorized
		case errors.Is(err, context.DeadlineExceeded):
			err = rpc.ErrTryAgainLater
		case errors.Is(err, context.Canceled):
			err = rpc.ErrTryAgainLater
		case errors.As(err, new(*mysql.MySQLError)):
			err = rpc.ErrInternal
		case errors.Is(err, repo.ErrSchemaVer):
			err = rpc.ErrInternal
		}
		return err
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
