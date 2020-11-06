package jsonrpc2

import (
	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/pkg/jsonrpc2x"
)

func (srv *Server) doExample(ctx Ctx, auth dom.Auth, arg api.RPCExampleReq, res *api.RPCExampleResp) error {
	userName, err := dom.ParseUserName(arg.UserName)
	if err != nil {
		return jsonrpc2x.ErrInvalidParams
	}

	example, err := srv.appl.Example(ctx, auth, *userName)
	if err == nil {
		*res = protoExample(*example)
	}
	return err
}

func (srv *Server) doIncExample(ctx Ctx, auth dom.Auth, _ api.RPCIncExampleReq, _ *api.RPCIncExampleResp) error {
	return srv.appl.IncExample(ctx, auth)
}
