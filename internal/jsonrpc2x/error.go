package jsonrpc2x

import (
	"errors"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

// Error wraps JSON-RPC 2.0 "Error object" to match (using errors.Is) any
// other JSON-RPC 2.0 error with same code.
type Error struct {
	Err *jsonrpc2.Error
}

// NewError returns an Error with given code and message.
func NewError(code int, message string) *Error {
	return &Error{Err: jsonrpc2.NewError(code, message)}
}

// Unwrap returns wrapped error.
func (e *Error) Unwrap() error {
	return e.Err
}

// Error returns JSON representation of Error.
func (e *Error) Error() string {
	return e.Err.Error()
}

// Is reports whether target error's code matches this error's code.
func (e *Error) Is(target error) bool {
	if rpcerr := new(jsonrpc2.Error); errors.As(target, &rpcerr) {
		return e.Err.Code == rpcerr.Code
	}
	return false
}
