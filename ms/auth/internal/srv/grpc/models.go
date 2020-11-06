package grpc

import (
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/powerman/go-monolith-example/api/proto/powerman/example/auth"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
)

func apiAccount(v app.User) *api.Account {
	return &api.Account{
		Name:       dom.NewName("accounts", v.Name.ID()).String(),
		User:       apiUser(v),
		Email:      v.Email,
		CreateTime: timestamppb.New(v.CreateTime),
	}
}

func apiUser(v app.User) *api.User {
	return &api.User{
		Name:        v.Name.String(),
		DisplayName: v.DisplayName,
		Access:      apiAccess(v),
	}
}

func apiAccess(v app.User) *api.Access {
	return &api.Access{
		Role: apiRole(v.Role),
	}
}

func apiRole(v app.Role) api.Access_Role {
	switch v {
	case app.RoleAdmin:
		return api.Access_ROLE_ADMIN
	case app.RoleUser:
		return api.Access_ROLE_USER
	default:
		panic(fmt.Sprintf("unknown app.Role: %v", v))
	}
}
