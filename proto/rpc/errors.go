// Package rpc defines JSON-RPC 2.0 errors used by all microservices.
package rpc

import "github.com/powerman/rpc-codec/jsonrpc2"

// All errors which may be returned by API methods.
//
//nolint:gomnd // By design.
var (
	ErrInvalidParams   = jsonrpc2.NewError(-32602, "invalid params")  // Client bug.
	ErrInternal        = jsonrpc2.NewError(-32000, "server error")    // Server bug or I/O issue.
	ErrTryAgainLater   = jsonrpc2.NewError(-503, "temporary error")   // Safe to resend.
	ErrTooManyRequests = jsonrpc2.NewError(-429, "too many requests") // Safe to resend (after delay).
	ErrNotFound        = jsonrpc2.NewError(-404, "not found")         // Given ID does not exists.
	ErrForbidden       = jsonrpc2.NewError(-403, "forbidden")         // Not allowed by permissions.
	ErrUnauthorized    = jsonrpc2.NewError(-401, "unauthorized")      // Missing or invalid Ctx.AccessToken.
)

// ErrsCommon may be returned by any API method.
var ErrsCommon = []error{
	ErrInvalidParams,
	ErrInternal,
	ErrTryAgainLater,
	ErrTooManyRequests,
	ErrForbidden,
	ErrUnauthorized,
}
