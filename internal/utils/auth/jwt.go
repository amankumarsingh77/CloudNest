package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTAuthenticator struct {
	Secret string
	Iss    string
	Aud    string
	Exp    time.Time
	NBF    time.Time
	IAT    time.Time
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		Secret: secret,
		Iss:    iss,
		Aud:    aud,
		Exp:    time.Now().Add(24 * time.Hour),
		NBF:    time.Now(),
		IAT:    time.Now(),
	}
}

func (s *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Secret))
}

func (s *JWTAuthenticator) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
