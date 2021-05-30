package dal

import (
	"time"
)

//nolint:gosec // False positive.
const (
	sqlUsersAdd = `
 INSERT INTO users (id, pass_salt, pass_hash, email, display_name, role)
 VALUES (:id, :pass_salt, :pass_hash, :email, :display_name, :role)
	`
	sqlUsersGet = `
 SELECT id, pass_salt, pass_hash, email, display_name, role, created_at
   FROM users
  WHERE id = :id
	`
	sqlUsersGetByEmail = `
 SELECT id, pass_salt, pass_hash, email, display_name, role, created_at
   FROM users
  WHERE LOWER(email) = LOWER(:email)
	`
	sqlGetUserByAccessToken = `
 SELECT id, email, display_name, role, u.created_at
   FROM access_tokens AS t LEFT JOIN users AS u ON (t.user_id = u.id)
  WHERE access_token = :access_token
	`
	sqlAccessTokensAdd = `
 INSERT INTO access_tokens (access_token, user_id)
 VALUES (:access_token, :user_id)
	`
	sqlAccessTokensDel = `
 DELETE FROM access_tokens
  WHERE access_token = :access_token
	`
	sqlAccessTokensDelByUser = `
 DELETE FROM access_tokens
  WHERE user_id = :user_id
	`
)

type (
	argUsersAdd struct {
		ID          string
		PassSalt    []byte
		PassHash    []byte
		Email       string
		DisplayName string
		Role        dbRole
	}

	argUsersGet struct {
		ID string
	}
	rowUsersGet struct {
		ID          string
		PassSalt    []byte
		PassHash    []byte
		Email       string
		DisplayName string
		Role        dbRole
		CreatedAt   time.Time
	}

	argUsersGetByEmail struct {
		Email string
	}
	rowUsersGetByEmail struct {
		ID          string
		PassSalt    []byte
		PassHash    []byte
		Email       string
		DisplayName string
		Role        dbRole
		CreatedAt   time.Time
	}

	argGetUserByAccessToken struct {
		AccessToken string
	}
	rowGetUserByAccessToken struct {
		ID          string
		Email       string
		DisplayName string
		Role        dbRole
		CreatedAt   time.Time
	}

	argAccessTokensAdd struct {
		AccessToken string
		UserID      string
	}

	argAccessTokensDel struct {
		AccessToken string
	}

	argAccessTokensDelByUser struct {
		UserID string
	}
)
