//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.$GOFILE Appl,Repo

// Package app provides business logic.
package app

import (
	"context"
	"errors"
	"time"

	"github.com/powerman/go-monolith-example/internal/dom"
)

// ServiceName provides name of this microservice for logs/metrics.
const ServiceName = "example"

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Errors.
var (
	ErrAccessDenied = errors.New("access denied")
	ErrNotFound     = errors.New("not found")
)

// Appl provides application features (use cases) service.
type Appl interface {
	// Example returns ...
	// Errors: ErrAccessDenied, ErrNotFound.
	Example(Ctx, dom.Auth, dom.UserID) (*Example, error)
	// IncExample creates or increments ...
	// Errors: none.
	IncExample(Ctx, dom.Auth) error
}

// Repo provides data storage.
type Repo interface {
	// Example returns ...
	// Errors: ErrNotFound.
	Example(Ctx, dom.UserID) (*Example, error)
	// IncExample creates or increments ...
	// Errors: none.
	IncExample(Ctx, dom.UserID) error
}

type (
	// Example describes ...
	Example struct {
		Counter int
		Mtime   time.Time
	}
)

// Config contains configuration for business-logic.
type Config struct{}

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
