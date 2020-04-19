// Package proto describes microservice example's public JSON-RPC 2.0 API.
package proto

import (
	"github.com/powerman/go-monolith-example/internal/dom"
)

type (
	// Example returns given user's Example.
	APIExampleReq struct {
		Ctx    `json:"Ctx"`
		UserID dom.UserID
	}
	APIExampleResp = Example

	// IncExample increments user's Example.
	APIIncExampleReq struct {
		Ctx `json:"Ctx"`
	}
	APIIncExampleResp struct{}
)
