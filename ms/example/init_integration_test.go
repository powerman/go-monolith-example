// +build integration

package example

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/gotest/testinit"
	"github.com/powerman/mysqlx"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/powerman/structlog"
	"github.com/spf13/cobra"
)

var rpcClient *jsonrpc2.Client

func init() { testinit.Setup(serialIntegration, setupIntegration) }

func setupIntegration() {
	const dir = "internal/migrations"
	const host = "localhost"
	log := structlog.FromContext(ctx, nil)

	cfgTempDB, cleanup, err := mysqlx.EnsureTempDB(log, "", def.TestMySQLCfg(def.MySQLAuth{
		User: def.ExampleDBUser,
		Pass: def.ExampleDBPass,
		DB:   def.ExampleDBName,
	}))
	if err != nil {
		testinit.Fatal(err)
	}
	testinit.Teardown(cleanup)

	def.RPCHost = host
	def.ExampleMetricsPort = 0
	def.ExampleRPCPort = unusedTCPPort(host)
	def.ExampleDBName = cfgTempDB.DBName
	def.ExampleGooseDir = dir

	serveCmd := &cobra.Command{}
	Service{}.Init(&cobra.Command{}, serveCmd)
	err = serveCmd.ParseFlags(nil)
	if err != nil {
		testinit.Fatal(err)
	}

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

	rpcAddr := fmt.Sprintf("%s:%d", host, def.ExampleRPCPort)
	waitTCPPort(rpcAddr)
	rpcClient = jsonrpc2.NewHTTPClient("http://" + rpcAddr + "/rpc")
}

func unusedTCPPort(host string) (port int) {
	var portStr string
	ln, err := net.Listen("tcp", host+":0")
	if err == nil {
		err = ln.Close()
	}
	if err == nil {
		_, portStr, err = net.SplitHostPort(ln.Addr().String())
	}
	if err == nil {
		port, err = strconv.Atoi(portStr)
	}
	if err != nil {
		testinit.Fatal(err)
	}
	return port
}

func waitTCPPort(addr string) {
	ctx, cancel := context.WithTimeout(ctx, def.TestSecond)
	defer cancel()
	var dialer net.Dialer
	for ; ctx.Err() != context.DeadlineExceeded; time.Sleep(def.TestSecond / 20) {
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			conn.Close()
			break
		}
	}
	if ctx.Err() != nil {
		testinit.Fatal("failed to start API")
	}
}

func call(method string, arg, res interface{}) error {
	return jsonrpc2.ServerError(rpcClient.Call(method, arg, res))
}
