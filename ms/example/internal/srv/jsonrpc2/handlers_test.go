package jsonrpc2_test

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/powerman/check"

	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/internal/apix"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/pkg/jsonrpc2x"
)

func TestExample(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, c, _, mockAppl := testNew(t)
	defer cleanup()

	exampleUser := &app.Example{Counter: 3}

	mockAppl.EXPECT().Example(gomock.Any(), authUser, authAdmin.UserID).Return(nil, app.ErrAccessDenied)
	mockAppl.EXPECT().Example(gomock.Any(), authAdmin, authUser.UserID).Return(exampleUser, nil)

	tests := []struct {
		token   apix.AccessToken
		userID  dom.UserID
		want    *api.Example
		wantErr error
	}{
		{tokenEmpty, authUser.UserID, nil, api.ErrUnauthorized},
		{tokenAdmin, userIDBad, nil, jsonrpc2x.ErrInvalidParams},
		{tokenUser, authAdmin.UserID, nil, api.ErrForbidden},
		{tokenAdmin, authUser.UserID, &api.Example{Counter: 3}, nil},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			req := api.RPCExampleReq{
				Ctx: apix.JSONRPC2Ctx{
					AccessToken: tc.token,
				},
				UserID: tc.userID,
			}
			var res api.RPCExampleResp
			err := c.Call("RPC.Example", req, &res)
			t.Err(err, tc.wantErr)
			if tc.wantErr == nil {
				t.DeepEqual(res, *tc.want)
			}
		})
	}
}

func TestIncExample(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	cleanup, c, _, mockAppl := testNew(t)
	defer cleanup()

	mockAppl.EXPECT().IncExample(gomock.Any(), authAdmin).Return(nil)

	req := api.RPCIncExampleReq{
		Ctx: apix.JSONRPC2Ctx{
			AccessToken: tokenAdmin,
		},
	}
	t.Nil(c.Call("RPC.IncExample", req, new(api.RPCIncExampleResp)))
}
