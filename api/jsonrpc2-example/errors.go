package api

import "github.com/powerman/rpc-codec/jsonrpc2"

// All generic errors which may be returned by RPC methods.
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

// ErrsCommon may be returned by any RPC method.
var ErrsCommon = []error{
	ErrInvalidParams,
	ErrInternal,
	ErrTryAgainLater,
	ErrTooManyRequests,
	ErrForbidden,
	ErrUnauthorized,
}

// All errors which may be returned by RPC methods.
//
//nolint:gomnd // By design.
var (
// ErrSomething = jsonrpc2.NewError(1, "Something is wrong")
)

// ErrsExtra list non-common errors which may be returned by concrete RPC method.
var ErrsExtra = map[string][]error{
	Name + ".Example":    {ErrNotFound},
	Name + ".IncExample": nil,
}
