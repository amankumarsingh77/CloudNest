package middleware

import (
	"github.com/amankumarsingh77/cloudnest/config"
	"net/http"
	"strconv"
)

func FileSizeLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit := config.Load().FileSizeLimit
		contentLength := r.Header.Get("Content-Length")
		if contentLength == "" {
			http.Error(w, "Missing Content-Length header", http.StatusBadRequest)
			return
		}

		fileSize, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Content-Length header", http.StatusBadRequest)
			return
		}

		if fileSize > limit {
			http.Error(w, "File size too large", http.StatusRequestEntityTooLarge)
			return
		}

		next.ServeHTTP(w, r)
	})
}
