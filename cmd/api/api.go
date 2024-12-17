package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/amankumarsingh77/cloudnest/internal/env"
	"github.com/amankumarsingh77/cloudnest/internal/handlers"
	middleware2 "github.com/amankumarsingh77/cloudnest/internal/middleware"
	"github.com/amankumarsingh77/cloudnest/internal/services"
	"github.com/amankumarsingh77/cloudnest/internal/store/db"
	"github.com/amankumarsingh77/cloudnest/internal/utils/auth"
	"github.com/amankumarsingh77/cloudnest/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	config     *env.Config
	dbStore    db.DbStore
	auth       auth.Authenticator
	middleware *middleware2.Middleware
	logger     *logger.Logger
	services   *services.Services
}

func (app *Application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.middleware.CORS)
	if app.config.Server.RateLimit.Enabled {
		r.Use(app.middleware.RateLimiter)
	}

	h := handlers.Handler{
		Auth:     app.auth,
		Services: app.services,
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", h.CreateUserHandler)
				r.Post("/signin", h.AuthenticateUserHandler)
				r.Group(func(r chi.Router) {
					r.Use(app.middleware.AuthToken)
					r.Patch("/", h.UpdateUserHandler)
				})
			})
			r.Route("/upload", func(r chi.Router) {
				r.Use(app.middleware.AuthToken)
				r.Post("/presignedurl", h.GetPresignedUrlHandler)
			})
			r.Route("/files", func(r chi.Router) {
				r.Use(app.middleware.AuthToken)
				r.Post("/", h.CreateFileHandler)
				r.Route("/{fileID}", func(r chi.Router) {
					r.Use(h.FileContextMiddleware)
					r.Get("/", h.GetFileHandler)
					r.Post("/", h.UpdateFileHandler)
				})
			})
		})
	})
	return r
}

func (app *Application) run(mux http.Handler) error {
	srv := http.Server{
		Addr:         app.config.Server.Port,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("Sever running on port %v", app.config.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}
	return nil
}
