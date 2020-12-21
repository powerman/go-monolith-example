// Package grpcx provide helpers for typical gRPC client/server.
package grpcx

import (
	"context"
	"time"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

const (
	keepaliveTime    = 50 * time.Second
	keepaliveTimeout = 10 * time.Second
	keepaliveMinTime = 30 * time.Second
	xForwardedFor    = "X-Forwarded-For"
)
