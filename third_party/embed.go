// Package third_party contains embedded files.
package third_party

import (
	"embed"
	"io/fs"
)

var (
	//go:embed swagger-ui/dist/favicon-16x16.png
	//go:embed swagger-ui/dist/favicon-32x32.png
	//go:embed swagger-ui/dist/oauth2-redirect.html
	//go:embed swagger-ui/dist/swagger-ui-bundle.js
	//go:embed swagger-ui/dist/swagger-ui.css
	//go:embed swagger-ui/dist/swagger-ui-standalone-preset.js
	swaggerUI embed.FS
	// SwaggerUI contains static files for Swagger UI required by
	// our ../web/static/swagger-ui/index.html.
	SwaggerUI, _ = fs.Sub(swaggerUI, "swagger-ui/dist")
)
