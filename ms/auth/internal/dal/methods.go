package dal

import (
	"database/sql"
	"errors"

	"github.com/powerman/sqlxx"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/pkg/repo"
)

func (r *Repo) AddUser(ctx Ctx, user app.User) error {
	return r.Tx(ctx, nil, func(tx *sqlxx.Tx) error {
		_, err := tx.NamedExecContext(ctx, sqlUsersAdd, argUsersAdd{
			ID:          user.Name.ID(),
			PassSalt:    user.PassHash.Salt,
			PassHash:    user.PassHash.Hash,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			Role:        dalRole(user.Role),
		})
		if repo.PostgresErrName(err, repo.PostgresUniqueViolation) {
			return app.ErrAlreadyExist
		}
		return err
	})
}

func (r *Repo) GetUser(ctx Ctx, userName dom.UserName) (res *app.User, err error) {
	err = r.Tx(ctx, &sql.TxOptions{ReadOnly: true}, func(tx *sqlxx.Tx) error {
		var resUsersGet rowUsersGet
		err := tx.NamedGetContext(ctx, &resUsersGet, sqlUsersGet, argUsersGet{
			ID: userName.ID(),
		})
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return app.ErrNotFound
		case err != nil:
			return err
		}
		res = appUserWithPass(resUsersGet)
		return nil
	})
	return
}

func (r *Repo) GetUserByEmail(ctx Ctx, email string) (res *app.User, err error) {
	err = r.Tx(ctx, &sql.TxOptions{ReadOnly: true}, func(tx *sqlxx.Tx) error {
		var resUsersGetByEmail rowUsersGetByEmail
		err := tx.NamedGetContext(ctx, &resUsersGetByEmail, sqlUsersGetByEmail, argUsersGetByEmail{
			Email: email,
		})
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return app.ErrNotFound
		case err != nil:
			return err
		}
		res = appUserWithPass(rowUsersGet(resUsersGetByEmail))
		return nil
	})
	return
}

func (r *Repo) GetUserByAccessToken(ctx Ctx, accessToken app.AccessToken) (res *app.User, err error) {
	err = r.Tx(ctx, &sql.TxOptions{ReadOnly: true}, func(tx *sqlxx.Tx) error {
		var resGetUserByAccessToken rowGetUserByAccessToken
		err := tx.NamedGetContext(ctx, &resGetUserByAccessToken, sqlGetUserByAccessToken, argGetUserByAccessToken{
			AccessToken: string(accessToken),
		})
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return app.ErrNotFound
		case err != nil:
			return err
		}
		res = appUser(resGetUserByAccessToken)
		return nil
	})
	return
}

func (r *Repo) AddAccessToken(ctx Ctx, accessToken app.AccessToken, userName dom.UserName) error {
	return r.Tx(ctx, nil, func(tx *sqlxx.Tx) error {
		_, err := tx.NamedExecContext(ctx, sqlAccessTokensAdd, argAccessTokensAdd{
			AccessToken: string(accessToken),
			UserID:      userName.ID(),
		})
		if repo.PostgresErrName(err, repo.PostgresForeignKeyViolation) {
			return app.ErrNotFound
		}
		return err
	})
}

func (r *Repo) DelAccessToken(ctx Ctx, accessToken app.AccessToken) error {
	return r.Tx(ctx, nil, func(tx *sqlxx.Tx) error {
		_, err := tx.NamedExecContext(ctx, sqlAccessTokensDel, argAccessTokensDel{
			AccessToken: string(accessToken),
		})
		return err
	})
}

func (r *Repo) DelAccessTokens(ctx Ctx, userName dom.UserName) error {
	return r.Tx(ctx, nil, func(tx *sqlxx.Tx) error {
		_, err := tx.NamedExecContext(ctx, sqlAccessTokensDelByUser, argAccessTokensDelByUser{
			UserID: userName.ID(),
		})
		return err
	})
}
