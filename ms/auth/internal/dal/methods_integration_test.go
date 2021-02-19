// build integration

package dal_test

import (
	"errors"
	"testing"
	"time"

	"github.com/powerman/check"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
)

var (
	errDupToken = errors.New(`duplicate key value violates unique constraint "access_tokens_pkey"`)
	now         = time.Now().Truncate(time.Second)
	tmplUAdmin  = app.User{
		Name: dom.NewUserName("admin"),
		PassHash: app.PassHash{
			Salt: []byte("saltAdmin"),
			Hash: []byte("hashAdmin"),
		},
		Email: "root@localhost",
		Role:  app.RoleAdmin,
	}
	tmplU1 = app.User{
		Name: dom.NewUserName("user1"),
		PassHash: app.PassHash{
			Salt: []byte("salt1"),
			Hash: []byte("hash1"),
		},
		Email:       "user1@localhost",
		DisplayName: "User 1",
		Role:        app.RoleUser,
	}
	tmplU2 = app.User{
		Name:  dom.NewUserName("user2"),
		Email: "user2@localhost",
		Role:  app.RoleUser,
	}
)

func TestUser(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	r := newTestRepo(t)

	var (
		uAdmin    = tmplUAdmin
		u1        = tmplU1
		u2        = tmplU2
		u1DupName = app.User{
			Name:  dom.NewUserName("user1"),
			Email: "user1dup@localhost",
			Role:  app.RoleUser,
		}
		u1DupEmail = app.User{
			Name:  dom.NewUserName("user1dup"),
			Email: "User1@LocalHost",
			Role:  app.RoleUser,
		}
	)

	res, err := r.GetUser(ctx, u1.Name)
	t.Err(err, app.ErrNotFound)
	t.Nil(res)
	res, err = r.GetUserByEmail(ctx, u1.Email)
	t.Err(err, app.ErrNotFound)
	t.Nil(res)

	tests := []struct {
		given   app.User
		wantErr error
	}{
		{uAdmin, nil},
		{uAdmin, app.ErrAlreadyExist},
		{u1, nil},
		{u1DupName, app.ErrAlreadyExist},
		{u1DupEmail, app.ErrAlreadyExist},
		{u2, nil},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			err := r.AddUser(ctx, tc.given)
			matchErr(t, err, tc.wantErr)
			if err == nil {
				if tc.given.PassHash.Salt == nil {
					tc.given.PassHash.Salt = []byte{}
					tc.given.PassHash.Hash = []byte{}
				}
				res, err := r.GetUser(ctx, tc.given.Name)
				t.Nil(err)
				t.GE(res.CreateTime, now)
				tc.given.CreateTime = res.CreateTime
				t.DeepEqual(res, &tc.given)

				res, err = r.GetUserByEmail(ctx, tc.given.Email)
				t.Nil(err)
				t.DeepEqual(res, &tc.given)
			}
		})
	}
}

func TestAccessToken(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	r := newTestRepo(t)

	var (
		uAdmin = tmplUAdmin
		u1     = tmplU1
		u2     = tmplU2
	)
	for _, u := range []*app.User{&uAdmin, &u1, &u2} {
		t.Nil(r.AddUser(ctx, *u))
		res, err := r.GetUser(ctx, u.Name)
		t.Nil(err)
		u.CreateTime = res.CreateTime
	}

	tests := []struct {
		AccessToken app.AccessToken
		userName    dom.UserName
		wantErr     error
	}{
		{"admintoken1", uAdmin.Name, nil},
		{"admintoken1", uAdmin.Name, errDupToken},
		{"admintoken2", uAdmin.Name, nil},
		{"u1token1", u1.Name, nil},
		{"u1token2", u1.Name, nil},
		{"u2token1", u2.Name, nil},
		{"u3token1", dom.NewUserName("user3"), app.ErrNotFound},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			err := r.AddAccessToken(ctx, tc.AccessToken, tc.userName)
			matchErr(t, err, tc.wantErr)
		})
	}

	t.Nil(r.DelAccessTokens(ctx, uAdmin.Name))
	t.Nil(r.DelAccessTokens(ctx, uAdmin.Name))
	t.Nil(r.DelAccessTokens(ctx, dom.NewUserName("nosuch")))
	t.Nil(r.DelAccessToken(ctx, "u1token1"))
	t.Nil(r.DelAccessToken(ctx, "u1token1"))
	t.Nil(r.DelAccessToken(ctx, "nosuch"))

	u1.PassHash = app.PassHash{}
	u2.PassHash = app.PassHash{}

	tests2 := []struct {
		AccessToken app.AccessToken
		want        *app.User
		wantErr     error
	}{
		{"nosuch", nil, app.ErrNotFound},
		{"admintoken1", nil, app.ErrNotFound},
		{"admintoken2", nil, app.ErrNotFound},
		{"u1token1", nil, app.ErrNotFound},
		{"u1token2", &u1, nil},
		{"u1token2", &u1, nil},
		{"u2token1", &u2, nil},
	}
	for _, tc := range tests2 {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := r.GetUserByAccessToken(ctx, tc.AccessToken)
			matchErr(t, err, tc.wantErr)
			t.DeepEqual(res, tc.want)
		})
	}
}
