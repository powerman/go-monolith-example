package grpc

import (
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/powerman/go-monolith-example/api/proto/powerman/example/auth"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
)

func apiAccount(m app.User) *api.Account {
	return &api.Account{
		Name:       dom.NewName("accounts", m.Name.ID()).String(),
		User:       apiUser(m),
		Email:      m.Email,
		CreateTime: timestamppb.New(m.CreateTime),
	}
}

func apiUser(m app.User) *api.User {
	return &api.User{
		Name:        m.Name.String(),
		DisplayName: m.DisplayName,
		Access:      apiAccess(m),
	}
}

func apiAccess(m app.User) *api.Access {
	return &api.Access{
		Role: apiRole(m.Role),
	}
}

func apiRole(m app.Role) api.Access_Role {
	switch m {
	case app.RoleAdmin:
		return api.Access_ROLE_ADMIN
	case app.RoleUser:
		return api.Access_ROLE_USER
	default:
		panic(fmt.Sprintf("unknown app.Role: %v", m))
	}
}
