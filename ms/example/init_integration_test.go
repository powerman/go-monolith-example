// +build integration

package example

import (
	"context"
	"fmt"

	"github.com/powerman/go-monolith-example/internal/apiauth"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/internal/netx"
	"github.com/powerman/gotest/testinit"
	"github.com/powerman/mysqlx"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/powerman/structlog"
)

var rpcClient *jsonrpc2.Client

func init() { testinit.Setup(serialIntegration, setupIntegration) }

func setupIntegration() {
	const rootDir = "../../"
	const host = "localhost"
	log := structlog.FromContext(ctx, nil)
	ctxStartup, cancel := context.WithTimeout(ctx, 3*def.TestSecond)
	defer cancel()

	cfgTempDB, cleanup, err := mysqlx.EnsureTempDB(log, "", cfg.MySQLConfig)
	if err != nil {
		testinit.Fatal(err)
	}
	testinit.Teardown(cleanup)

	cfg.MySQLConfig = cfgTempDB
	cfg.GooseDir = rootDir + cfg.GooseDir
	cfg.MetricsAddr = netx.NewAddr(host, 0)
	cfg.RPCAddr = netx.NewAddr(host, netx.UnusedTCPPort(host))

	authn = mockAuthn{}

	ctxShutdown, shutdown := context.WithCancel(ctx)
	errc := make(chan error)
	go func() { errc <- Service{}.Serve(ctx, ctxShutdown, shutdown) }()
	testinit.Teardown(func() {
		shutdown()
		if err := <-errc; err != nil {
			log.Println("failed to Serve:", err)
		}
		repo.Close()
	})

	if netx.WaitTCPPort(ctxStartup, cfg.RPCAddr) != nil {
		testinit.Fatal("failed to connect to API")
	}
	rpcClient = jsonrpc2.NewHTTPClient(fmt.Sprintf("http://%s/rpc", cfg.RPCAddr))
}

func call(method string, arg, res interface{}) error {
	return jsonrpc2.ServerError(rpcClient.Call(method, arg, res))
}

type mockAuthn struct{}

func (mockAuthn) Authenticate(token apiauth.AccessToken) (dom.Auth, error) {
	switch token {
	case tokenAdmin:
		return authAdmin, nil
	case tokenUser:
		return authUser, nil
	default:
		return dom.Auth{}, apiauth.ErrAccessTokenInvalid
	}
}
