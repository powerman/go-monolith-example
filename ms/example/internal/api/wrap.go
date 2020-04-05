//go:generate genny -in=$GOFILE -out=gen.$GOFILE gen "Example=IncExample"

package api

import (
	"github.com/powerman/go-monolith-example/internal/jsonrpc2x"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/proto/rpc"
)

// Example implements JSON-RPC 2.0 method.
func (api *API) Example(arg rpc.ExampleReq, res *rpc.ExampleResp) error {
	ctx, log, methodName, auth, err := arg.NewContext(app.ServiceName)
	recovery := jsonrpc2x.MakeRecovery(log, app.Metric)
	metrics := jsonrpc2x.MakeMetrics(metric, methodName)
	apiError := makeAPIError(rpc.ErrsExtra[methodName])
	accessLog := jsonrpc2x.MakeAccessLog(log)
	handler := recovery(metrics(apiError(accessLog(func() error {
		if err != nil {
			return err
		}
		return api.doExample(ctx, auth, arg, res)
	}))))
	return handler()
}
