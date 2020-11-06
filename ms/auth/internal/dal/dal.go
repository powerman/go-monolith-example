package dal

import (
	"context"
	"sync"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
)

type Ctx = context.Context

type Repo struct {
	sync.Mutex
	users  map[dom.UserName]app.User
	tokens map[app.AccessToken]dom.UserName
}

func New() *Repo {
	return &Repo{
		users:  make(map[dom.UserName]app.User),
		tokens: make(map[app.AccessToken]dom.UserName),
	}
}

func (r *Repo) AddUser(ctx Ctx, user app.User) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.users[user.Name]; ok {
		return app.ErrAlreadyExist
	}
	for i := range r.users {
		if user.Email == r.users[i].Email {
			return app.ErrAlreadyExist
		}
	}
	r.users[user.Name] = user
	return nil
}

func (r *Repo) GetUser(ctx Ctx, userName dom.UserName) (*app.User, error) {
	r.Lock()
	defer r.Unlock()
	user, ok := r.users[userName]
	if !ok {
		return nil, app.ErrNotFound
	}
	return &user, nil
}

func (r *Repo) GetUserByEmail(ctx Ctx, email string) (*app.User, error) {
	r.Lock()
	defer r.Unlock()
	for _, user := range r.users { //nolint:gocritic // rangeValCopy.
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, app.ErrNotFound
}

func (r *Repo) GetUserByAccessToken(ctx Ctx, accessToken app.AccessToken) (*app.User, error) {
	r.Lock()
	defer r.Unlock()
	userName, ok := r.tokens[accessToken]
	if !ok {
		return nil, app.ErrNotFound
	}
	user, ok := r.users[userName]
	if !ok {
		return nil, app.ErrNotFound
	}
	return &user, nil
}

func (r *Repo) AddAccessToken(ctx Ctx, userName dom.UserName) (app.AccessToken, error) {
	r.Lock()
	defer r.Unlock()
	accessToken := app.AccessToken(dom.NewID())
	r.tokens[accessToken] = userName // XXX May overwrite existing record.
	return accessToken, nil
}

func (r *Repo) DelAccessToken(ctx Ctx, accessToken app.AccessToken) error {
	r.Lock()
	defer r.Unlock()
	delete(r.tokens, accessToken)
	return nil
}

func (r *Repo) DelAccessTokens(ctx Ctx, userName dom.UserName) error {
	r.Lock()
	defer r.Unlock()
	for accessToken, name := range r.tokens { //nolint:gocritic // rangeValCopy.
		if name.String() == userName.String() {
			delete(r.tokens, accessToken)
		}
	}
	return nil
}
