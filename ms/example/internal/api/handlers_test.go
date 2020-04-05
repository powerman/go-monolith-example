package api

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/proto"
	"github.com/powerman/go-monolith-example/proto/rpc"
)

func TestExample(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	api, mockAppl := testNew(t)

	var (
		exampleUser = &app.Example{Counter: 3}
	)

	mockAppl.EXPECT().Example(gomock.Any(), authUser, authAdmin.UserID).Return(nil, app.ErrAccessDenied)
	mockAppl.EXPECT().Example(gomock.Any(), authAdmin, authUser.UserID).Return(exampleUser, nil)

	tests := []struct {
		token   proto.AccessToken
		userID  dom.UserID
		want    *app.Example
		wantErr error
	}{
		{tokenEmpty, authUser.UserID, nil, rpc.ErrUnauthorized},
		{tokenAdmin, userIDBad, nil, rpc.ErrInvalidParams},
		{tokenUser, authAdmin.UserID, nil, rpc.ErrForbidden},
		{tokenAdmin, authUser.UserID, exampleUser, nil},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			req := rpc.ExampleReq{
				Ctx: rpc.Ctx{
					AccessToken: tc.token,
				},
				UserID: tc.userID,
			}
			var res rpc.ExampleResp
			err := api.Example(req, &res)
			t.Err(err, tc.wantErr)
			if tc.wantErr == nil {
				t.DeepEqual(res, proto.Example(*tc.want))
			}
		})
	}
}

func TestIncExample(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	api, mockAppl := testNew(t)

	mockAppl.EXPECT().IncExample(gomock.Any(), authAdmin).Return(nil)

	req := rpc.IncExampleReq{
		Ctx: rpc.Ctx{
			AccessToken: tokenAdmin,
		},
	}
	t.Nil(api.IncExample(req, new(rpc.IncExampleResp)))
}
