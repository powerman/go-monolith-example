// Package proto describes entities used by public API.
package proto

import "time"

type (
	// AccessToken should be either "admin" or "user".
	AccessToken string

	// Example is an example of user's data.
	Example struct {
		Counter int
		Mtime   time.Time `json:"updated_at"`
	}
)
