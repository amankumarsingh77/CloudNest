package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

type ClientRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewClientRateLimiter(requestPerSec int, burst int) *ClientRateLimiter {
	return &ClientRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(float64(requestPerSec) / 60),
		burst:    burst,
	}
}

func (c *ClientRateLimiter) GetLimiter(ip string) *rate.Limiter {
	c.mu.Lock()
	defer c.mu.Unlock()
	limiter, ok := c.limiters[ip]
	if !ok {
		limiter = rate.NewLimiter(c.rate, c.burst)
		c.limiters[ip] = limiter
	}
	return limiter
}

func RateLimiterMiddleware(requestPerSec int, burst int) func(handler http.Handler) http.Handler {
	rateLimiter := NewClientRateLimiter(requestPerSec, burst)

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			limiter := rateLimiter.GetLimiter(ip)
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}
