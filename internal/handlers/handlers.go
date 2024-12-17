package handlers

import (
	"github.com/amankumarsingh77/cloudnest/internal/services"
	"github.com/amankumarsingh77/cloudnest/internal/utils/auth"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Handler struct {
	Auth     auth.Authenticator
	Services *services.Services
}
