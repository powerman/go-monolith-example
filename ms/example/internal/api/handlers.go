package api

import (
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/proto"
	"github.com/powerman/go-monolith-example/proto/rpc"
)

func (api *API) doExample(ctx Ctx, auth dom.Auth, arg rpc.ExampleReq, res *rpc.ExampleResp) error {
	if arg.UserID <= 0 {
		return rpc.ErrInvalidParams
	}

	example, err := api.a.Example(ctx, auth, arg.UserID)
	if err == nil {
		*res = proto.Example(*example)
	}
	return err
}

func (api *API) doIncExample(ctx Ctx, auth dom.Auth, _ rpc.IncExampleReq, _ *rpc.IncExampleResp) error {
	return api.a.IncExample(ctx, auth)
}
