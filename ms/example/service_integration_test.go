// +build integration

package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/mysqlx"

	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/internal/apix"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/jsonrpc2x"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

func TestSmoke(tt *testing.T) {
	t := check.T(tt)
	ctrl := gomock.NewController(t)

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
	t.Must(t.Nil(netx.WaitTCPPort(ctxStartup, cfg.BindAddr), "connect to service"))

	mockAuthn.EXPECT().Authenticate(gomock.Any(), apix.AccessToken(tokenAdmin)).Return(authAdmin, nil).AnyTimes()
	mockAuthn.EXPECT().Authenticate(gomock.Any(), apix.AccessToken(tokenUser)).Return(authUser, nil).AnyTimes()
	mockAuthn.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(dom.Auth{}, apix.ErrAccessTokenInvalid).AnyTimes()

	rpcClient := jsonrpc2x.NewHTTPClient(fmt.Sprintf("http://%s/rpc", cfg.BindAddr))

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
		argExample.UserName = authAdmin.UserName.String()
		err := rpcClient.Call("RPC.Example", argExample, &resExample)
		t.Err(err, api.ErrNotFound)
	}
	{
		argExample.Ctx.AccessToken = tokenUser
		argExample.UserName = authUser.UserName.String()
		t.Nil(rpcClient.Call("RPC.Example", argExample, &resExample))
		t.NotZero(resExample.UpdatedAt)
		t.DeepEqual(resExample, api.RPCExampleResp{
			Counter:   2,
			UpdatedAt: resExample.UpdatedAt,
		})
	}
}
