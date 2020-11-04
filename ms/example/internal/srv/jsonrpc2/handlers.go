package jsonrpc2

import (
	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/pkg/jsonrpc2x"
)

func (srv *Server) doExample(ctx Ctx, auth dom.Auth, arg api.RPCExampleReq, res *api.RPCExampleResp) error {
	if arg.UserID <= 0 {
		return jsonrpc2x.ErrInvalidParams
	}

	example, err := srv.appl.Example(ctx, auth, arg.UserID)
	if err == nil {
		*res = protoExample(*example)
	}
	return err
}

func (srv *Server) doIncExample(ctx Ctx, auth dom.Auth, _ api.RPCIncExampleReq, _ *api.RPCIncExampleResp) error {
	return srv.appl.IncExample(ctx, auth)
}
