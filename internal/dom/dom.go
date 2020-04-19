// Package dom contains common domain (business-logic) entities.
package dom

// UserID is a user ID.
type UserID int

// Auth should contain all authentication and authorization info
// needed to execute any operation on behalf of some user.
type Auth struct {
	UserID UserID
	Admin  bool
}
