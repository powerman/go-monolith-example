// +build integration

package example

import (
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/go-monolith-example/proto/rpc"
)

func TestSmoke(tt *testing.T) {
	t := check.T(tt)
	var (
		argIncExample rpc.IncExampleReq
		resIncExample rpc.IncExampleResp
		argExample    rpc.ExampleReq
		resExample    rpc.ExampleResp
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
		argExample.UserID = userIDAdmin
		err := call("API.Example", argExample, &resExample)
		t.Err(err, rpc.ErrNotFound)
	}
	{
		argExample.Ctx.AccessToken = tokenUser
		argExample.UserID = userIDUser
		t.Nil(call("API.Example", argExample, &resExample))
		t.NotZero(resExample.Mtime)
		t.DeepEqual(resExample, rpc.ExampleResp{
			Counter: 2,
			Mtime:   resExample.Mtime,
		})
	}
}
