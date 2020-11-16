package app

import (
	"fmt"
	"regexp"
	"time"

	"github.com/powerman/sensitive"

	"github.com/powerman/go-monolith-example/internal/dom"
)

var reValidUserID = regexp.MustCompile(`^[a-z0-9-]{4,63}$`)

func (a *App) Register(ctx Ctx, userID string, password sensitive.String, user *User) error {
	if userID == "" {
		userID = dom.NewID()
	}
	if !reValidUserID.MatchString(userID) {
		return fmt.Errorf("%w: userID should be 4-63 characters [a-z0-9-]", ErrValidate)
	}
	user.Name = dom.NewUserName(userID)
	user.PassHash = a.newPassHash(password, nil)
	user.Role = RoleUser
	if user.Name.ID() == "admin" {
		user.Role = RoleAdmin
	}
	user.CreateTime = time.Now()
	return a.repo.AddUser(ctx, *user)
}

func (a *App) LoginByUserID(ctx Ctx, userID string, password sensitive.String) (AccessToken, error) {
	userName := dom.NewUserName(userID)
	user, err := a.repo.GetUser(ctx, userName)
	if err != nil {
		return "", err
	}
	if !a.equalPassHash(password, user.PassHash) {
		return "", ErrWrongPassword
	}
	accessToken := AccessToken(dom.NewID())
	err = a.repo.AddAccessToken(ctx, accessToken, user.Name)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (a *App) LoginByEmail(ctx Ctx, email string, password sensitive.String) (AccessToken, error) {
	user, err := a.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	return a.LoginByUserID(ctx, user.Name.ID(), password)
}

func (a *App) Authenticate(ctx Ctx, accessToken AccessToken) (*User, error) {
	return a.repo.GetUserByAccessToken(ctx, accessToken)
}

func (a *App) Logout(ctx Ctx, accessToken AccessToken) error {
	return a.repo.DelAccessToken(ctx, accessToken)
}

func (a *App) LogoutUser(ctx Ctx, userName dom.UserName) error {
	return a.repo.DelAccessTokens(ctx, userName)
}
