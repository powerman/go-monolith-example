//go:generate gobin -m -run github.com/golang/mock/mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.$GOFILE Authn

package apix

import (
	"errors"
	"fmt"

	"github.com/powerman/sensitive"

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
	Authenticate(Ctx, AccessToken) (dom.Auth, error)
}

// AccessToken is an access token.
type AccessToken string

// Format wraps sensitive.String.
func (s AccessToken) Format(f fmt.State, c rune) { sensitive.String(s).Format(f, c) }

// MarshalJSON wraps sensitive.String.
func (s AccessToken) MarshalJSON() ([]byte, error) { return sensitive.String(s).MarshalJSON() }

// MarshalText wraps sensitive.String.
func (s AccessToken) MarshalText() ([]byte, error) { return sensitive.String(s).MarshalText() }
