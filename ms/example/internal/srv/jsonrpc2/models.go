package jsonrpc2

import (
	api "github.com/powerman/go-monolith-example/api/jsonrpc2-example"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
)

func protoExample(v app.Example) api.Example {
	return api.Example{
		Counter:   v.Counter,
		UpdatedAt: v.Mtime,
	}
}
