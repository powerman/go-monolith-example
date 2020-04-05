package rpc

import "github.com/powerman/rpc-codec/jsonrpc2"

// All errors which may be returned by API methods.
var (
	ErrInvalidParams   = jsonrpc2.NewError(-32602, "invalid params")  // client bug
	ErrInternal        = jsonrpc2.NewError(-32000, "server error")    // server bug or I/O issue
	ErrTryAgainLater   = jsonrpc2.NewError(-503, "temporary error")   // safe to resend
	ErrTooManyRequests = jsonrpc2.NewError(-429, "too many requests") // safe to resend (after delay)
	ErrNotFound        = jsonrpc2.NewError(-404, "not found")         // given ID not exists
	ErrForbidden       = jsonrpc2.NewError(-403, "forbidden")         // not allowed by permissions
	ErrUnauthorized    = jsonrpc2.NewError(-401, "unauthorized")      // no/wrong Ctx.AccessToken
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

// ErrsExtra list non-common errors which may be returned by concrete API method.
var ErrsExtra = map[string][]error{
	"Example":    {ErrNotFound},
	"IncExample": nil,
}
