package api

import (
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/proto/rpc"
	proto "github.com/powerman/go-monolith-example/proto/rpc-example"
)

func (api *API) doExample(ctx Ctx, auth dom.Auth, arg proto.APIExampleReq, res *proto.APIExampleResp) error {
	if arg.UserID <= 0 {
		return rpc.ErrInvalidParams
	}

	example, err := api.a.Example(ctx, auth, arg.UserID)
	if err == nil {
		*res = protoExample(*example)
	}
	return err
}

func (api *API) doIncExample(ctx Ctx, auth dom.Auth, _ proto.APIIncExampleReq, _ *proto.APIIncExampleResp) error {
	return api.a.IncExample(ctx, auth)
}
