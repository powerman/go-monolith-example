package api

import (
	"time"

	"github.com/powerman/go-monolith-example/internal/apix"
)

// Ctx must be embedded and tagged `json:"Ctx"` to prevent JSON embedding.
type Ctx = apix.JSONRPC2Ctx

type (
	// Example is an example of user's data.
	Example struct {
		Counter   int
		UpdatedAt time.Time
	}
)
