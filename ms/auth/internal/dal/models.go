package dal

import (
	"fmt"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
)

type dbRole string

const (
	roleAdmin dbRole = "admin"
	roleUser  dbRole = "user"
)

func dalRole(role app.Role) dbRole {
	switch role {
	case app.RoleAdmin:
		return roleAdmin
	case app.RoleUser:
		return roleUser
	default:
		panic(fmt.Sprintf("unknown app.Role: %v", role))
	}
}

func appRole(role dbRole) app.Role {
	switch role {
	case roleAdmin:
		return app.RoleAdmin
	case roleUser:
		return app.RoleUser
	default:
		panic(fmt.Sprintf("unknown dal.Role: %v", role))
	}
}

func appUserWithPass(row rowUsersGet) *app.User {
	return &app.User{
		Name: dom.NewUserName(row.ID),
		PassHash: app.PassHash{
			Salt: row.PassSalt,
			Hash: row.PassHash,
		},
		Email:       row.Email,
		DisplayName: row.DisplayName,
		Role:        appRole(row.Role),
		CreateTime:  row.CreatedAt,
	}
}

func appUser(row rowGetUserByAccessToken) *app.User {
	return &app.User{
		Name:        dom.NewUserName(row.ID),
		Email:       row.Email,
		DisplayName: row.DisplayName,
		Role:        appRole(row.Role),
		CreateTime:  row.CreatedAt,
	}
}
