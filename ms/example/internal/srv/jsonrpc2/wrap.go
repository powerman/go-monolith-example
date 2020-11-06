//go:generate gobin -m -run github.com/cheekybits/genny -in=$GOFILE -out=gen.$GOFILE gen "Example=IncExample"
//go:generate sed -i -e "\\,^//go:generate,d" gen.$GOFILE

package jsonrpc2

import (
	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/pkg/jsonrpc2x"
)

// Example implements JSON-RPC 2.0 method.
func (srv *Server) Example(arg api.RPCExampleReq, res *api.RPCExampleResp) error {
	ctx, log, methodName, auth, err := arg.NewContext(srv.authn, app.ServiceName)
	validateErr := jsonrpc2x.MakeValidateErr(log, srv.cfg.StrictErr, api.ErrsCommon, api.ErrsExtra["RPC."+methodName])
	recovery := jsonrpc2x.MakeRecovery(log, app.Metric)
	metrics := jsonrpc2x.MakeMetrics(metric, "RPC."+methodName)
	accessLog := jsonrpc2x.MakeAccessLog(log)
	handler := validateErr(recovery(metrics(apiErr(accessLog(func() error {
		if err != nil {
			return err
		}
		if auth.UserName == dom.NoUser {
			return api.ErrUnauthorized
		}
		return srv.doExample(ctx, auth, arg, res)
	})))))
	return handler()
}
