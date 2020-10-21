// +build integration

package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/mysqlx"
	"github.com/powerman/rpc-codec/jsonrpc2"

	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/internal/apix"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

func TestSmoke(tt *testing.T) {
	t := check.T(tt)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthn := apix.NewMockAuthn(ctrl)

	s := &Service{cfg: cfg}
	s.authn = mockAuthn

	tempDBCfg, cleanup, err := mysqlx.EnsureTempDB(tLogger(*t), "", cfg.MySQL)
	cfg.MySQL = tempDBCfg // Assign to cfg and not s.cfg as a reminder: they are the same.
	t.Must(t.Nil(err))
	defer cleanup()

	ctxStartup, cancel := context.WithTimeout(ctx, def.TestTimeout)
	defer cancel()
	ctxShutdown, shutdown := context.WithCancel(ctx)
	errc := make(chan error)
	go func() { errc <- s.RunServe(ctxStartup, ctxShutdown, shutdown) }()
	defer func() {
		shutdown()
		t.Nil(<-errc, "RunServe")
		if s.repo != nil {
			s.repo.Close()
		}
	}()
	t.Must(t.Nil(netx.WaitTCPPort(ctxStartup, cfg.Addr), "connect to service"))

	mockAuthn.EXPECT().Authenticate(tokenAdmin).Return(authAdmin, nil).AnyTimes()
	mockAuthn.EXPECT().Authenticate(tokenUser).Return(authUser, nil).AnyTimes()
	mockAuthn.EXPECT().Authenticate(gomock.Any()).Return(dom.Auth{}, apix.ErrAccessTokenInvalid).AnyTimes()

	rpcClient := jsonrpc2.NewHTTPClient(fmt.Sprintf("http://%s/rpc", cfg.Addr))

	var (
		argIncExample api.RPCIncExampleReq
		resIncExample api.RPCIncExampleResp
		argExample    api.RPCExampleReq
		resExample    api.RPCExampleResp
	)

	{ // insert
		argIncExample.Ctx.AccessToken = tokenUser
		t.Nil(rpcClient.Call("RPC.IncExample", argIncExample, &resIncExample))
	}
	{ // update
		t.Nil(rpcClient.Call("RPC.IncExample", argIncExample, &resIncExample))
	}
	{
		argExample.Ctx.AccessToken = tokenAdmin
		argExample.UserID = authAdmin.UserID
		err := rpcClient.Call("RPC.Example", argExample, &resExample)
		t.Err(jsonrpc2.ServerError(err), api.ErrNotFound)
	}
	{
		argExample.Ctx.AccessToken = tokenUser
		argExample.UserID = authUser.UserID
		t.Nil(rpcClient.Call("RPC.Example", argExample, &resExample))
		t.NotZero(resExample.UpdatedAt)
		t.DeepEqual(resExample, api.RPCExampleResp{
			Counter:   2,
			UpdatedAt: resExample.UpdatedAt,
		})
	}
}
