package proto

import (
	"time"

	"github.com/powerman/go-monolith-example/internal/apiauth"
)

// Ctx must be embedded and tagged `json:"Ctx"` to prevent JSON embedding.
type Ctx = apiauth.Ctx

type (
	// Example is an example of user's data.
	Example struct {
		Counter   int
		UpdatedAt time.Time
	}
)
