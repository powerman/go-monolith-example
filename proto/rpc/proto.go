// Package rpc describes public JSON-RPC 2.0 API.
package rpc

import (
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/proto"
)

// Tag `json:"Ctx"` is required to prevent JSON embedding.

type (
	// Example returns given user's Example.
	ExampleReq struct {
		Ctx    `json:"Ctx"`
		UserID dom.UserID
	}
	ExampleResp = proto.Example

	// IncExample increments user's Example.
	IncExampleReq struct {
		Ctx `json:"Ctx"`
	}
	IncExampleResp struct{}
)
