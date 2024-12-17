package db

import (
	"context"
	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
	"github.com/amankumarsingh77/cloudnest/internal/services/s3"
	"github.com/amankumarsingh77/cloudnest/internal/store/db"
)

type DBService struct {
	File interface {
		CreateFile(ctx context.Context, file *entities.File) error
		UpdateFile(ctx context.Context, file *entities.File) error
		GetFileById(ctx context.Context, id string) (*entities.File, error)
	}
	Folder interface {
		CreateFolder(ctx context.Context, folder *entities.Folder) error
		DeleteFolder(context context.Context, id string) error
		UpdateFolder(ctx context.Context, folder *entities.Folder) error
		GetFilesInFolder(context context.Context, id string) (*[]entities.File, error)
	}
	User interface {
		CreateUser(ctx context.Context, user *entities.User) error
		UpdateUser(ctx context.Context, user *entities.User) error
		GetUserById(ctx context.Context, userId string) (*entities.User, error)
		GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
		GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	}
	UserQuota interface {
		GetUserQuota(ctx context.Context, userId string) (*entities.UserQuota, error)
		CreateUserQuota(ctx context.Context, userQuota *entities.UserQuota) error
		UpdateUserQuota(ctx context.Context, userQuota *entities.UserQuota) error
	}
}

func NewDbService(db *db.DbStore, service *s3.S3Service) DBService {
	return DBService{
		File:      &FileService{db, service},
		Folder:    &FolderService{db},
		User:      &UserService{db},
		UserQuota: &UserQuotaService{db},
	}
}
