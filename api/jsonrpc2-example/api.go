// Package api describes microservice example's public JSON-RPC 2.0 API.
//
//nolint:tagliatelle // Pascal case instead of camel.
package api

// Name is a net/rpc type name used as a prefix before method names.
const Name = "RPC"

type (
	// RPC.Example returns given user's Example.
	RPCExampleReq struct {
		Ctx      `json:"Ctx"`
		UserName string
	}
	RPCExampleResp = Example

	// RPC.IncExample increments user's Example.
	RPCIncExampleReq struct {
		Ctx `json:"Ctx"`
	}
	RPCIncExampleResp struct{}
)
