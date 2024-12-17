package middleware

import (
	"context"
	"github.com/amankumarsingh77/cloudnest/internal/services"
	"github.com/amankumarsingh77/cloudnest/internal/utils/auth"
	"github.com/amankumarsingh77/cloudnest/internal/utils/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func AuthTokenMiddleware(auth *auth.JWTAuthenticator, service *services.Services) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				json.WriteJsonError(w, http.StatusUnauthorized, "Authorization header is missing")
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 && parts[0] != "Bearer" {
				json.WriteJsonError(w, http.StatusUnauthorized, "Authorization header is invalid")
				return
			}
			tokenString := parts[1]
			jwtToken, err := auth.ValidateToken(tokenString)
			if err != nil {
				json.WriteJsonError(w, http.StatusUnauthorized, "Authorization header is invalid")
				return
			}
			claims := jwtToken.Claims.(jwt.MapClaims)
			userId := claims["sub"].(string)
			ctx := r.Context()
			user, err := service.DB.User.GetUserById(ctx, userId)
			if err != nil {
				json.WriteJsonError(w, http.StatusUnauthorized, "Authorization header is invalid")
				return
			}
			ctx = context.WithValue(ctx, "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}
