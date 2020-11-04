// Package dom contains common domain (business-logic) entities.
package dom

import (
	"fmt"
	"strings"
)

// Auth should contain all authentication and authorization info
// needed to execute any operation on behalf of some user.
type Auth struct {
	UserName UserName
	Admin    bool
}

const (
	pfxUsers = "users/"
)

// UserName is a user name.
type UserName string

// Empty values.
const (
	NoUser UserName = ""
)

// NewUserName converts user's ID to Name.
func NewUserName(userID string) UserName {
	if userID == "" {
		panic("no userID")
	}
	return UserName(pfxUsers + userID)
}

// Valid returns true if name is valid by format (but it may not exists).
func (name UserName) Valid() bool {
	return len(name) > len(pfxUsers) && strings.HasPrefix(string(name), pfxUsers)
}

// ID converts user's Name to ID.
func (name UserName) ID() string {
	id := strings.TrimPrefix(string(name), pfxUsers)
	if len(id) == len(name) {
		panic(fmt.Sprintf("invalid UserName: %q", name))
	}
	return id
}
