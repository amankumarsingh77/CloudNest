package middleware

import (
	"github.com/amankumarsingh77/cloudnest/internal/services"
	"github.com/amankumarsingh77/cloudnest/internal/utils/auth"
	"net/http"
)

type Middleware struct {
	AuthToken   func(next http.Handler) http.Handler
	CORS        func(handler http.Handler) http.Handler
	RateLimiter func(handler http.Handler) http.Handler
}

func NewMiddleware(auth *auth.JWTAuthenticator, service *services.Services, allowedDomains []string, requestsPerSec, burst int) *Middleware {

	return &Middleware{
		AuthToken:   AuthTokenMiddleware(auth, service),
		CORS:        CorsMiddleware(allowedDomains),
		RateLimiter: RateLimiterMiddleware(requestsPerSec, burst),
	}
}
