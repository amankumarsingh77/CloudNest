package middleware

import (
	"net/http"

	"github.com/amankumarsingh77/cloudnest/internal/utils/json"
)

func CorsMiddleware(allowedDomains []string) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			var isAllowed bool
			for _, allowedOrigin := range allowedDomains {
				if allowedOrigin == origin {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				json.WriteJsonError(w, http.StatusForbidden, "Access denied : cors error")
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE,PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Length, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}
