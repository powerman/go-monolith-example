// Package api describes microservice example's public JSON-RPC 2.0 API.
package api

import (
	"github.com/powerman/go-monolith-example/internal/dom"
)

type (
	// RPC.Example returns given user's Example.
	RPCExampleReq struct {
		Ctx    `json:"Ctx"`
		UserID dom.UserID
	}
	RPCExampleResp = Example

	// RPC.IncExample increments user's Example.
	RPCIncExampleReq struct {
		Ctx `json:"Ctx"`
	}
	RPCIncExampleResp struct{}
)
