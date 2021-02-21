// Package proto contains embedded files.
package proto

import "embed"

// OpenAPI contains path/to/*.swagger.json for public gRPC services.
//go:embed */*/*/*.swagger.json
var OpenAPI embed.FS
