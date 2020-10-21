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
	cleanup, a, mockRepo := testNew(t)
	defer cleanup()

	exampleUser := &app.Example{Counter: 3}

	mockRepo.EXPECT().Example(gomock.Any(), authUser.UserID).Return(exampleUser, nil).Times(2)
	mockRepo.EXPECT().Example(gomock.Any(), gomock.Any()).Return(nil, app.ErrNotFound)

	tests := []struct {
		auth    dom.Auth
		userID  dom.UserID
		want    *app.Example
		wantErr error
	}{
		{authUser, authUser.UserID, exampleUser, nil},
		{authUser, authAdmin.UserID, nil, app.ErrAccessDenied},
		{authAdmin, authUser.UserID, exampleUser, nil},
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
	cleanup, a, mockRepo := testNew(t)
	defer cleanup()

	mockRepo.EXPECT().IncExample(gomock.Any(), authAdmin.UserID).Return(nil)
	t.Nil(a.IncExample(ctx, authAdmin))
}
