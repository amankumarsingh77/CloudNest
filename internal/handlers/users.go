package handlers

import (
	"context"
	"errors"
	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
	"github.com/amankumarsingh77/cloudnest/internal/store/db"
	"github.com/amankumarsingh77/cloudnest/internal/utils/json"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"time"
)

type createUserPayload struct {
	Username string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required,max=20"`
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required,max=30"`
}

type updateUserPayload struct {
	Username *string `json:"username" validate:"omitempty,max=20"`
	Email    *string `json:"email" validate:"omitempty,max=100"`
	Name     *string `json:"name" validate:"omitempty,max=30"`
}

type authenticateUserPayload struct {
	Username string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required,max=20"`
}

type UserWithToken struct {
	*entities.User
	Token string `json:"token"`
}

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userPayload createUserPayload
	if err := json.ReadJson(w, r, &userPayload); err != nil {
		json.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(&userPayload); err != nil {
		json.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	user := entities.User{
		Name:     userPayload.Name,
		Email:    userPayload.Email,
		Username: userPayload.Username,
	}
	err := user.Password.Set(userPayload.Password)
	if err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ctx := r.Context()
	//TODO validation pending
	if err := h.Services.DB.User.CreateUser(ctx, &user); err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	slog.Info("user registered successfully ")
	if err := json.WriteJson(w, http.StatusCreated, user); err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) AuthenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userPayload authenticateUserPayload
	if err := json.ReadJson(w, r, &userPayload); err != nil {
		json.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(&userPayload); err != nil {
		json.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()
	user, err := h.Services.DB.User.GetUserByUsername(ctx, userPayload.Username)
	if err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err = user.Password.Check(userPayload.Password); err != nil {
		json.WriteJsonError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"iss": "cloudnest",
		"aud": "cloudnest",
	}

	token, err := h.Auth.GenerateToken(claims)
	if err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	userWithToken := UserWithToken{
		user,
		token,
	}
	if err = json.WriteJson(w, http.StatusOK, userWithToken); err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r.Context())
	var userUpdatePayload updateUserPayload
	if err := json.ReadJson(w, r, &userUpdatePayload); err != nil {
		json.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(&userUpdatePayload); err != nil {
		json.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	if userUpdatePayload.Username != nil {
		user.Username = *userUpdatePayload.Username
	}
	if userUpdatePayload.Email != nil {
		user.Email = *userUpdatePayload.Email
	}
	if userUpdatePayload.Name != nil {
		user.Name = *userUpdatePayload.Name
	}
	ctx := r.Context()
	if err := h.Services.DB.User.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			json.WriteJsonError(w, http.StatusNotFound, err.Error())
			return
		} else {
			json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err := json.WriteJson(w, http.StatusOK, user); err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) UserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "userID")
		if id == "" {
			json.WriteJsonError(w, http.StatusBadRequest, "id required")
			return
		}
		ctx := r.Context()
		user, err := h.Services.DB.User.GetUserById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, db.ErrorNotFound):
				json.WriteJsonError(w, http.StatusNotFound, err.Error())
			default:
				json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(ctx context.Context) *entities.User {
	return ctx.Value("user").(*entities.User)
}
