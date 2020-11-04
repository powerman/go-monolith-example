package app

import "github.com/powerman/go-monolith-example/internal/dom"

// Example implements Appl interface.
func (a *App) Example(ctx Ctx, auth dom.Auth, userID dom.UserName) (*Example, error) {
	if !(auth.UserName == userID || auth.Admin) {
		metric.ErrAccessDeniedTotal.Inc()
		return nil, ErrAccessDenied
	}
	return a.repo.Example(ctx, userID)
}

// IncExample implements Appl interface.
func (a *App) IncExample(ctx Ctx, auth dom.Auth) error {
	return a.repo.IncExample(ctx, auth.UserName)
}
