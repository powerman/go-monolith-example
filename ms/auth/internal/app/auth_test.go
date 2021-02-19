package app_test

import (
	"io"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/sensitive"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
)

func TestRegister(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	a, mockRepo := testNew(t)

	var (
		now   = time.Now()
		uAnon = &app.User{
			Role: app.RoleUser,
		}
		uUser = &app.User{
			Name:        dom.NewUserName("user"),
			DisplayName: "UseR",
			Role:        app.RoleUser,
		}
		uAdmin = &app.User{
			Name:  dom.NewUserName("admin"),
			Email: "root@host",
			Role:  app.RoleAdmin,
		}
	)

	mockRepo.EXPECT().AddUser(gomock.Any(), app.MatchUser{uAnon})
	mockRepo.EXPECT().AddUser(gomock.Any(), app.MatchUser{&app.User{Role: app.RoleUser}}).Return(app.ErrAlreadyExist)
	mockRepo.EXPECT().AddUser(gomock.Any(), app.MatchUser{uUser})
	mockRepo.EXPECT().AddUser(gomock.Any(), app.MatchUser{uAdmin})

	tests := []struct {
		userID   string
		password string
		user     *app.User
		want     *app.User
		wantErr  error
	}{
		{"", "", &app.User{}, uAnon, nil},
		{"", "", &app.User{}, nil, app.ErrAlreadyExist},
		{"bad", "", &app.User{}, nil, app.ErrValidate},
		{"user", "pass", &app.User{DisplayName: "UseR"}, uUser, nil},
		{"admin", "secret", &app.User{Email: "root@host"}, uAdmin, nil},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			err := a.Register(ctx, tc.userID, sensitive.String(tc.password), tc.user)
			t.Err(err, tc.wantErr)
			if err == nil {
				t.DeepEqual(tc.user, tc.want)
				t.Greater(tc.user.CreateTime, now)
			}
		})
	}
}

func TestLoginByUserID(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	a, mockRepo := testNew(t)

	mockRepo.EXPECT().AddUser(gomock.Any(), gomock.Any())
	uAdmin := &app.User{}
	t.Nil(a.Register(ctx, "admin", "secret", uAdmin))
	mockRepo.EXPECT().GetUser(gomock.Any(), uAdmin.Name).Return(uAdmin, nil).AnyTimes()
	mockRepo.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound)
	mockRepo.EXPECT().AddAccessToken(gomock.Any(), gomock.Any(), uAdmin.Name).Return(nil)
	mockRepo.EXPECT().AddAccessToken(gomock.Any(), gomock.Any(), uAdmin.Name).Return(io.EOF)

	tests := []struct {
		userID  string
		pass    string
		wantErr error
	}{
		{"user", "", app.ErrNotFound},
		{"admin", "wrong", app.ErrWrongPassword},
		{"admin", "secret", nil},
		{"admin", "secret", io.EOF},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := a.LoginByUserID(ctx, tc.userID, sensitive.String(tc.pass))
			t.Err(err, tc.wantErr)
			if tc.wantErr == nil {
				t.Len(res, 26)
			} else {
				t.Len(res, 0)
			}
		})
	}
}

func TestLoginByEmail(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	a, mockRepo := testNew(t)

	mockRepo.EXPECT().AddUser(gomock.Any(), gomock.Any())
	uAdmin := &app.User{Email: "admin@host"}
	t.Nil(a.Register(ctx, "admin", "secret", uAdmin))
	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), uAdmin.Email).Return(uAdmin, nil).AnyTimes()
	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound)
	mockRepo.EXPECT().GetUser(gomock.Any(), uAdmin.Name).Return(uAdmin, nil).AnyTimes()
	mockRepo.EXPECT().AddAccessToken(gomock.Any(), gomock.Any(), uAdmin.Name).Return(nil)
	mockRepo.EXPECT().AddAccessToken(gomock.Any(), gomock.Any(), uAdmin.Name).Return(io.EOF)

	tests := []struct {
		email   string
		pass    string
		wantErr error
	}{
		{"user@host", "", app.ErrNotFound},
		{"admin@host", "wrong", app.ErrWrongPassword},
		{"admin@host", "secret", nil},
		{"admin@host", "secret", io.EOF},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := a.LoginByEmail(ctx, tc.email, sensitive.String(tc.pass))
			t.Err(err, tc.wantErr)
			if tc.wantErr == nil {
				t.Len(res, 26)
			} else {
				t.Len(res, 0)
			}
		})
	}
}
