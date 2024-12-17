package services

import (
	"context"
	"github.com/amankumarsingh77/cloudnest/internal/env"
	dbservice "github.com/amankumarsingh77/cloudnest/internal/services/db"
	"github.com/amankumarsingh77/cloudnest/internal/services/s3"
	"github.com/amankumarsingh77/cloudnest/internal/store/db"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type Services struct {
	Storage StorageService
	DB      dbservice.DBService
	// Add other services here as needed
}

// StorageService interface defines storage operations
type StorageService interface {
	GetPresignedURL(ctx context.Context, filename string, filesize int64) (string, error)
	Delete(ctx context.Context, filename string) error
	Get(ctx context.Context, filename string) (string, error)
}

func NewServices(store *db.DbStore, s3Client *awsS3.Client, cfg *env.Config) *Services {
	storageService := s3.NewS3Service(s3Client, awsS3.NewPresignClient(s3Client), &cfg.Storage)
	dbService := dbservice.NewDbService(store, storageService)
	return &Services{
		Storage: storageService,
		DB:      dbService,
	}
}
