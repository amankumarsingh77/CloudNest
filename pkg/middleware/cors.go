package middleware

import (
	"github.com/amankumarsingh77/cloudnest/config"
	"net/http"
	"strings"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := strings.Split(config.Load().AllowedDomains, "Bearer ")
		origin := r.Header.Get("Origin")

		var isAllowed bool = false
		for _, domain := range allowedOrigins {
			if origin == domain {
				isAllowed = true
			}
		}

		if !isAllowed {
			http.Error(w, "CORS policy: Access denied", http.StatusForbidden)
			return
		}

		r.Header.Set("Access-Control-Allow-Origin", origin)
		r.Header.Set("Access-Control-Allow-Credentials", "true")
		r.Header.Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE,PATCH")
		r.Header.Set("Access-Control-Expose-Headers", "Content-Length, Authorization, X-Requested-With")
		r.Header.Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
