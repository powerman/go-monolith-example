//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.$GOFILE Authenticator

package apiauth

import (
	"errors"

	"github.com/powerman/go-monolith-example/internal/dom"
)

// Errors.
var (
	ErrAccessTokenInvalid = errors.New("invalid AccessToken")
)

// Authenticator validates AccessToken.
type Authenticator interface {
	// Authenticate validates AccessToken and returns corresponding
	// Auth. If validation fails returns zero Auth with error.
	//
	// Errors: ErrAccessTokenInvalid.
	Authenticate(AccessToken) (dom.Auth, error)
}

// AccessToken is an access token.
type AccessToken string
