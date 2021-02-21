// Package web contains embedded files.
package web

import (
	"embed"
	"io/fs"
)

var (
	//go:embed static/swagger-ui/index.html
	swaggerUI embed.FS
	// SwaggerUI contains overlay files for third_party.SwaggerUI.
	SwaggerUI, _ = fs.Sub(swaggerUI, "static/swagger-ui")
)
