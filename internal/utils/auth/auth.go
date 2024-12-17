package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	GenerateToken(clams jwt.Claims) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type Password struct {
	Text *string
	Hash []byte
}

func (p *Password) Set(s string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Text = &s
	p.Hash = hash
	return nil
}

func (p *Password) Check(password string) error {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(password))
}
