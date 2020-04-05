package app

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/go-monolith-example/internal/dom"
)

func TestExample(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	a, mockRepo := testNew(t)

	var (
		exampleUser = &Example{Counter: 3}
	)

	mockRepo.EXPECT().Example(gomock.Any(), authUser.UserID).Return(exampleUser, nil).Times(2)
	mockRepo.EXPECT().Example(gomock.Any(), gomock.Any()).Return(nil, ErrNotFound)

	tests := []struct {
		auth    dom.Auth
		userID  dom.UserID
		want    *Example
		wantErr error
	}{
		{authUser, authUser.UserID, exampleUser, nil},
		{authUser, authAdmin.UserID, nil, ErrAccessDenied},
		{authAdmin, authUser.UserID, exampleUser, nil},
		{authAdmin, userIDBad, nil, ErrNotFound},
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

	mockRepo.EXPECT().IncExample(gomock.Any(), authAdmin.UserID).Return(nil)
	t.Nil(a.IncExample(ctx, authAdmin))
}
