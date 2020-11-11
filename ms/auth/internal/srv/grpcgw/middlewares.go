package grpcgw

import (
	"net/http"

	"github.com/rs/cors"
)

func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", "0")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		next.ServeHTTP(w, r)
	})
}

func corsAllowAll(next http.Handler) http.Handler {
	return cors.AllowAll().Handler(next)
}
