//go:generate genny -in=$GOFILE -out=gen.$GOFILE gen "Example=IncExample"

package api

import (
	"github.com/powerman/go-monolith-example/internal/jsonrpc2x"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	proto "github.com/powerman/go-monolith-example/proto/rpc-example"
)

// Example implements JSON-RPC 2.0 method.
func (api *API) Example(arg proto.APIExampleReq, res *proto.APIExampleResp) error {
	ctx, log, methodName, auth, err := arg.NewContext(api.authn, app.ServiceName)
	validateErr := jsonrpc2x.MakeValidateErr(log, api.strictErr, proto.ErrsExtra[methodName])
	recovery := jsonrpc2x.MakeRecovery(log, app.Metric)
	metrics := jsonrpc2x.MakeMetrics(metric, methodName)
	accessLog := jsonrpc2x.MakeAccessLog(log)
	handler := validateErr(recovery(metrics(protoErr(jsonrpc2x.ProtoErr(accessLog(func() error {
		if err != nil {
			return err
		}
		return api.doExample(ctx, auth, arg, res)
	}))))))
	return handler()
}
