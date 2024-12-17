package main

import (
	"log"

	"github.com/amankumarsingh77/cloudnest/internal/env"
	"github.com/amankumarsingh77/cloudnest/internal/middleware"
	"github.com/amankumarsingh77/cloudnest/internal/services"
	"github.com/amankumarsingh77/cloudnest/internal/services/s3"
	"github.com/amankumarsingh77/cloudnest/internal/store/db"
	"github.com/amankumarsingh77/cloudnest/internal/utils/auth"
	apiLogger "github.com/amankumarsingh77/cloudnest/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Panic(err)
	}

	cfg := env.NewConfig()
	err := cfg.Validate()
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	dbCon, err := db.New(
		cfg.Database.Url,
		cfg.Database.MaxIdleCons,
		cfg.Database.MaxOpenCons,
		cfg.Database.MaxIdleTime,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbCon.Close()

	storage := db.NewDbStore(dbCon)
	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.Security.JWTSecret, cfg.Security.AUD, cfg.Security.AUD)
	logger := apiLogger.NewLogger()

	s3Client, err := s3.NewS3Client(
		cfg.Storage.S3Region,
		cfg.Storage.S3AccessKey,
		cfg.Storage.S3SecretKey,
	)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	appServices := services.NewServices(&storage, s3Client, cfg)

	middlewares := middleware.NewMiddleware(
		jwtAuthenticator,
		appServices,
		cfg.Security.AllowedOrigins,
		cfg.Server.RateLimit.RequestsPerMinute,
		cfg.Server.RateLimit.BurstSize,
	)
	app := Application{
		config:     cfg,
		dbStore:    storage,
		auth:       jwtAuthenticator,
		middleware: middlewares,
		logger:     logger,
		services:   appServices,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
