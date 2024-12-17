package entities

import "time"

type Claims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Token struct {
	AccessToken string `json:"access_token"`
}
