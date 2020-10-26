//go:generate gobin -m -run github.com/golang/mock/mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.$GOFILE Authn

package apix

import (
	"errors"

	"github.com/powerman/go-monolith-example/internal/dom"
)

// Errors.
var (
	ErrAccessTokenInvalid = errors.New("invalid AccessToken")
)

// Authn validates AccessToken.
type Authn interface {
	// Authenticate validates AccessToken and returns corresponding
	// Auth. If validation fails returns zero Auth with error.
	//
	// Errors: ErrAccessTokenInvalid.
	Authenticate(AccessToken) (dom.Auth, error)
}

// AccessToken is an access token.
type AccessToken string
