package serve

import (
	"crypto/tls"
	"net/http"
	"net/rpc"

	"github.com/powerman/rpc-codec/jsonrpc2"

	"github.com/powerman/go-monolith-example/pkg/netx"
)

// RPC starts HTTP server on addr path /rpc using rcvr as JSON-RPC 2.0
// handler.
func RPC(ctx Ctx, addr netx.Addr, tlsConfig *tls.Config, rcvr interface{}) error {
	return RPCName(ctx, addr, tlsConfig, rcvr, "")
}

// RPCName starts HTTP server on addr path /rpc using rcvr as JSON-RPC 2.0
// handler but uses the provided name for the type instead of the
// receiver's concrete type.
func RPCName(ctx Ctx, addr netx.Addr, tlsConfig *tls.Config, rcvr interface{}, name string) (err error) {
	srv := rpc.NewServer()
	if name != "" {
		err = srv.RegisterName(name, rcvr)
	} else {
		err = srv.Register(rcvr)
	}
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/rpc", jsonrpc2.HTTPHandler(srv))
	return HTTP(ctx, addr, tlsConfig, mux, "JSON-RPC 2.0")
}
