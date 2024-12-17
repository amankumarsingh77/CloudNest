package entities

import (
	"github.com/amankumarsingh77/cloudnest/internal/utils/auth"
	"github.com/amankumarsingh77/cloudnest/internal/utils/json"
)

type User struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	Username      string          `json:"username"`
	Password      auth.Password   `json:"-"`
	Status        string          `json:"status"`
	EmailVerified bool            `json:"email_verified"`
	LastLogin     json.NullString `json:"last_login_at"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}
