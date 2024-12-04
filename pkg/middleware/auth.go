package middleware

import (
	"context"
	"github.com/amankumarsingh77/cloudnest/config"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is empty", http.StatusUnauthorized)
			return
		}
		tokenString := strings.Split(authHeader, " ")[1]
		if tokenString == "" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Load().JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
