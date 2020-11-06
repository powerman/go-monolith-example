package apix

import "github.com/powerman/go-monolith-example/internal/dom"

// AccessTokenParser handles example tokens.
type AccessTokenParser map[AccessToken]dom.Auth

// NewAccessTokenParser creates AccessTokenParser which validates JWT
// signature using rsaPublicKey in PEM format.
func NewAccessTokenParser() AccessTokenParser {
	return AccessTokenParser{
		"admin": dom.Auth{UserName: dom.NewUserName("1"), Admin: true},
		"user":  dom.Auth{UserName: dom.NewUserName("2")}, //nolint:gomnd // Example.
	}
}

// Authenticate implements Authenticator interface.
func (p AccessTokenParser) Authenticate(_ Ctx, token AccessToken) (auth dom.Auth, err error) {
	auth = p[token]
	if auth.UserName == dom.NoUser {
		err = ErrAccessTokenInvalid
	}
	return auth, err
}
