// Package dom contains common domain (business-logic) entities.
package dom

import (
	"crypto/rand"
	"strings"

	"github.com/oklog/ulid"
)

// Auth should contain all authentication and authorization info
// needed to execute any operation on behalf of some user.
type Auth struct {
	UserName UserName
	Admin    bool
}

// NewID returns unique ID with up to 63 [a-z0-9] characters.
func NewID() string {
	return strings.ToLower(ulid.MustNew(ulid.Now(), rand.Reader).String())
}
