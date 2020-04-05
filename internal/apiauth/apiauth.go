// Package apiauth provide helpers for API to check authentication.
package apiauth

import (
	"errors"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/proto"
)

// Errors.
var (
	ErrInvalidAccessToken = errors.New("invalid AccessToken")
)

// ParseAccessToken validates access token and returns corresponding auth.
func ParseAccessToken(token proto.AccessToken) (auth dom.Auth, err error) {
	switch token {
	case "admin":
		auth = dom.Auth{
			UserID: 1,
			Admin:  true,
		}
	case "user":
		auth = dom.Auth{
			UserID: 2, //nolint:gomnd // Just an example.
		}
	default:
		err = ErrInvalidAccessToken
	}
	return auth, err
}
