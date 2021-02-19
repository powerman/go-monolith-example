package app_test

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/powerman/check"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
)

func TestExample(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	a, mockRepo := testNew(t)

	exampleUser := &app.Example{Counter: 3}

	mockRepo.EXPECT().Example(gomock.Any(), authUser.UserName).Return(exampleUser, nil).Times(2)
	mockRepo.EXPECT().Example(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound)

	tests := []struct {
		auth    dom.Auth
		userID  dom.UserName
		want    *app.Example
		wantErr error
	}{
		{authUser, authUser.UserName, exampleUser, nil},
		{authUser, authAdmin.UserName, nil, app.ErrAccessDenied},
		{authAdmin, authUser.UserName, exampleUser, nil},
		{authAdmin, userIDBad, nil, app.ErrNotFound},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := a.Example(ctx, tc.auth, tc.userID)
			t.Err(err, tc.wantErr)
			t.Equal(res, tc.want)
		})
	}
}

func TestIncExample(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	a, mockRepo := testNew(t)

	mockRepo.EXPECT().IncExample(gomock.Any(), authAdmin.UserName).Return(nil)
	t.Nil(a.IncExample(ctx, authAdmin))
}
