package api

import (
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	proto "github.com/powerman/go-monolith-example/proto/rpc-example"
)

func protoExample(v app.Example) proto.Example {
	return proto.Example{
		Counter:   v.Counter,
		UpdatedAt: v.Mtime,
	}
}
