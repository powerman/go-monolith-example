//go:generate -command mockgen sh -c "$(git rev-parse --show-toplevel)/.gobincache/$DOLLAR{DOLLAR}0 \"$DOLLAR{DOLLAR}@\"" mockgen
//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.$GOFILE

// Package app provides business logic.
package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/powerman/sensitive"

	"github.com/powerman/go-monolith-example/internal/dom"
)

// ServiceName provides name of this microservice for logs/metrics.
const ServiceName = "auth"

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Errors.
var (
	ErrAccessDenied  = errors.New("access denied")
	ErrAlreadyExist  = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
	ErrValidate      = errors.New("validate")
	ErrWrongPassword = errors.New("wrong password")
)

// Appl provides application features (use cases) service.
type Appl interface {
	// Register creates and returns new user account.
	// User can provide optional userID (username).
	// These fields will be ignored in input and set automatically:
	// Name, PassHash, Role, CreateTime.
	// If userID=="admin" then user's role will be set to RoleAdmin.
	// Errors: ErrAlreadyExist, ErrValidate.
	Register(_ Ctx, userID string, password sensitive.String, _ *User) error
	// LoginByUserID creates and returns AccessToken for the user.
	// Errors: ErrNotFound, ErrWrongPassword.
	LoginByUserID(_ Ctx, userID string, password sensitive.String) (AccessToken, error)
	// LoginByEmail work in same way as LoginByUserID.
	LoginByEmail(_ Ctx, email string, password sensitive.String) (AccessToken, error)
	// Authenticate returns identity tied to AccessToken.
	// Errors: ErrNotFound.
	Authenticate(Ctx, AccessToken) (*User, error)
	// Logout invalidates given AccessToken.
	// Errors: none.
	Logout(_ Ctx, _ AccessToken) error
	// LogoutUser invalidates all user's AccessToken.
	// Errors: none.
	LogoutUser(_ Ctx, _ dom.UserName) error
}

// Repo provides data storage.
type Repo interface {
	// AddUser creates user.
	// Errors: ErrAlreadyExist.
	AddUser(Ctx, User) error
	// GetUser reads User by UserName.
	// Errors: ErrNotFound.
	GetUser(Ctx, dom.UserName) (*User, error)
	// GetUserByEmail reads User by email.
	// Errors: ErrNotFound.
	GetUserByEmail(Ctx, string) (*User, error)
	// GetUserByAccessToken reads User by AccessToken.
	// Errors: ErrNotFound.
	GetUserByAccessToken(Ctx, AccessToken) (*User, error)
	// AddAccessToken creates and returns AccessToken for given user.
	// Errors: ErrNotFound.
	AddAccessToken(Ctx, AccessToken, dom.UserName) error
	// DelAccessToken deletes given AccessToken.
	// Errors: none.
	DelAccessToken(Ctx, AccessToken) error
	// DelAccessTokens deletes all AccessToken for given user.
	// Errors: none.
	DelAccessTokens(Ctx, dom.UserName) error
}

type (
	// User contains data needed for authentication, identity and
	// permissions.
	User struct {
		Name        dom.UserName
		PassHash    PassHash
		Email       string
		DisplayName string
		Role        Role
		CreateTime  time.Time
	}
	// PassHash contains hashed password.
	PassHash struct {
		Salt sensitive.Bytes
		Hash sensitive.Bytes
	}
	// Role defines possible roles for a user.
	Role int
	// AccessToken is a token tied to some identity and permissions.
	AccessToken string
)

// Format wraps sensitive.String.
func (s AccessToken) Format(f fmt.State, c rune) { sensitive.String(s).Format(f, c) }

// MarshalJSON wraps sensitive.String.
func (s AccessToken) MarshalJSON() ([]byte, error) { return sensitive.String(s).MarshalJSON() }

// MarshalText wraps sensitive.String.
func (s AccessToken) MarshalText() ([]byte, error) { return sensitive.String(s).MarshalText() }

// Roles.
//go:generate -command stringer sh -c "$(git rev-parse --show-toplevel)/.gobincache/$DOLLAR{DOLLAR}0 \"$DOLLAR{DOLLAR}@\"" stringer
//go:generate stringer -output=stringer.Role.go -type=Role -trimprefix=Role
const (
	_ Role = iota
	RoleAdmin
	RoleUser
	roleMax
)

// Config contains configuration for business-logic.
type Config struct {
	Secret sensitive.Bytes
}

// App implements interface Appl.
type App struct {
	cfg  Config
	repo Repo
}

// New creates and returns new App.
func New(repo Repo, cfg Config) *App {
	a := &App{
		cfg:  cfg,
		repo: repo,
	}
	return a
}
