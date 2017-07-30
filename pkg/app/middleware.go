package app

import (
	"net/http"

	"github.com/acoshift/middleware"
)

// OnlyRootPath serve handler only path == "/"
func OnlyRootPath() middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				NotFound(w, r)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
