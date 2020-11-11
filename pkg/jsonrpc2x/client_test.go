package jsonrpc2x_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/rpc-codec/jsonrpc2"

	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/jsonrpc2x"
	"github.com/powerman/go-monolith-example/pkg/netx"
	"github.com/powerman/go-monolith-example/pkg/serve"
)

func TestCall(tt *testing.T) {
	t := check.T(tt)
	endpoint, noBody := testNewServerProxy(t)
	c := jsonrpc2x.NewHTTPClient(endpoint)
	err42 := jsonrpc2.NewError(-32601, "rpc: can't find service no.such")

	t.Err(c.Call("no.such", nil, nil), err42)
	noBody()
	t.Err(c.Call("no.such", nil, nil), io.ErrUnexpectedEOF)
	t.Err(c.Call("no.such", nil, nil), err42)
}

func testNewServerProxy(t *check.C) (endpoint string, noBody func()) {
	endpoint = testNewServer(t)

	backend, err := url.Parse(endpoint)
	t.Must(t.Nil(err))
	backend.Path = ""
	proxy := httputil.NewSingleHostReverseProxy(backend)
	frontend := httptest.NewServer(proxy)
	t.Cleanup(frontend.Close)
	endpoint = frontend.URL + "/rpc"

	noBody = func() {
		proxy.ModifyResponse = func(w *http.Response) error {
			w.Body = ioutil.NopCloser(bytes.NewReader(nil))
			proxy.ModifyResponse = nil
			return nil
		}
	}

	return endpoint, noBody
}

func testNewServer(t *check.C) string {
	rpcAddr := netx.NewAddr("localhost", netx.UnusedTCPPort("localhost"))
	ctx, cancel := context.WithCancel(context.Background())
	errc := make(chan error, 1)
	go func() { errc <- serve.RPC(ctx, rpcAddr, nil, &TestRPC{}) }()

	t.Cleanup(func() {
		t.Helper()
		cancel()
		t.Nil(<-errc, "serve.RPC")
	})

	ctxStartup, cancelStartup := context.WithTimeout(ctx, def.TestTimeout)
	defer cancelStartup()
	t.Must(t.Nil(netx.WaitTCPPort(ctxStartup, rpcAddr), "connect to service"))
	endpoint := fmt.Sprintf("http://%s/rpc", rpcAddr)
	return endpoint
}

type TestRPC struct{}

func (TestRPC) Method(_ struct{}, _ *struct{}) error { return nil }
