// +build integration

package example

import (
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/go-monolith-example/proto/rpc"
	proto "github.com/powerman/go-monolith-example/proto/rpc-example"
)

func TestSmoke(tt *testing.T) {
	t := check.T(tt)
	var (
		argIncExample proto.APIIncExampleReq
		resIncExample proto.APIIncExampleResp
		argExample    proto.APIExampleReq
		resExample    proto.APIExampleResp
	)

	{ // insert
		argIncExample.Ctx.AccessToken = tokenUser
		t.Nil(call("API.IncExample", argIncExample, &resIncExample))
	}
	{ // update
		t.Nil(call("API.IncExample", argIncExample, &resIncExample))
	}
	{
		argExample.Ctx.AccessToken = tokenAdmin
		argExample.UserID = authAdmin.UserID
		err := call("API.Example", argExample, &resExample)
		t.Err(err, rpc.ErrNotFound)
	}
	{
		argExample.Ctx.AccessToken = tokenUser
		argExample.UserID = authUser.UserID
		t.Nil(call("API.Example", argExample, &resExample))
		t.NotZero(resExample.UpdatedAt)
		t.DeepEqual(resExample, proto.APIExampleResp{
			Counter:   2,
			UpdatedAt: resExample.UpdatedAt,
		})
	}
}
