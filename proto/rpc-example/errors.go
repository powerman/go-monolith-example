package proto

import "github.com/powerman/go-monolith-example/proto/rpc"

// ErrsExtra list non-common errors which may be returned by concrete API method.
var ErrsExtra = map[string][]error{
	"Example":    {rpc.ErrNotFound},
	"IncExample": nil,
}
