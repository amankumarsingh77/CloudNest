package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
)

type DbStore struct {
	Users interface {
		Create(context.Context, *entities.User) error
		Update(context.Context, *entities.User) error
		GetUserById(context.Context, string) (*entities.User, error)
		GetUserByEmail(context.Context, string) (*entities.User, error)
		GetUserByUsername(context.Context, string) (*entities.User, error)
		GetUserQuota(context.Context, string) (*entities.UserQuota, error)
	}
	Files interface {
		//Create(context.Context, *sql.Tx, *entities.File) error
		//Update(context.Context, *sql.Tx, *entities.File) error
		//Delete(context.Context, string) error
		GetFileById(context.Context, string) (*entities.File, error)
		CreateAndUpdateQuotaAndVersion(context.Context, *entities.File) error
		UpdateAndUpdateQuotaAndVersion(context.Context, *entities.File) error
	}
	Folders interface {
		Create(context.Context, *entities.Folder) error
		Delete(context.Context, string) error
		GetFilesInFolder(context.Context, string) (*[]entities.File, error)
	}
}

type StorageError struct {
	Code    string
	Message string
	Err     error
}

func (e *StorageError) Error() string {
	return e.Message
}

var (
	ErrorNotFound        = &StorageError{Code: "NOT_FOUND", Message: "Resource not found"}
	ErrDuplicate         = &StorageError{Code: "DUPLICATE", Message: "Resource already exists"}
	ErrQuotaExceeded     = &StorageError{Code: "QUOTA_EXCEEDED", Message: "quota exceeded"}
	QueryTimeOutDuration = time.Second * 5
)

func NewDbStore(db *sql.DB) DbStore {
	return DbStore{
		Users:   &UserStore{db},
		Files:   &FileStore{db},
		Folders: &FolderStore{db},
		//UserQuota: &UserQuotaStore{db},
	}
}
