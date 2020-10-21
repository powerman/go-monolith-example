//go:generate genny -in=$GOFILE -out=gen.$GOFILE gen "Example=IncExample"

package jsonrpc2

import (
	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/internal/jsonrpc2x"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
)

// Example implements JSON-RPC 2.0 method.
func (srv *Server) Example(arg api.RPCExampleReq, res *api.RPCExampleResp) error {
	ctx, log, methodName, auth, err := arg.NewContext(srv.authn, app.ServiceName)
	validateErr := jsonrpc2x.MakeValidateErr(log, srv.cfg.StrictErr, api.ErrsExtra["RPC."+methodName])
	recovery := jsonrpc2x.MakeRecovery(log, app.Metric)
	metrics := jsonrpc2x.MakeMetrics(metric, "RPC."+methodName)
	accessLog := jsonrpc2x.MakeAccessLog(log)
	handler := validateErr(recovery(metrics(jsonrpc2x.APIErr(apiErr(accessLog(func() error {
		if err != nil {
			return err
		}
		if auth.UserID == dom.NoUserID {
			return api.ErrUnauthorized
		}
		return srv.doExample(ctx, auth, arg, res)
	}))))))
	return handler()
}
