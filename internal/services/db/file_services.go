package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
	"github.com/amankumarsingh77/cloudnest/internal/services/s3"
	"github.com/amankumarsingh77/cloudnest/internal/store/db"
)

type FileService struct {
	db        *db.DbStore
	s3Service *s3.S3Service
}

func (s *FileService) CreateFile(ctx context.Context, file *entities.File) error {
	userQuota, err := s.db.Users.GetUserQuota(ctx, file.CreatedBy)
	if err != nil {
		return err
	}
	if !userQuota.IsAllowedToUpload {
		return errors.New("quota exceeded : please upgrade")
	}
	remQuota := userQuota.StorageLimit - userQuota.StorageUsed
	if remQuota < file.Size {
		return errors.New("file size exceed quota limit")
	}
	// Check the file integrity with s3
	if !s.s3Service.FileExist(ctx, file.RemoteFileName) {
		return errors.New("file integrity check failed")
	}

	if err = s.db.Files.CreateAndUpdateQuotaAndVersion(ctx, file); err != nil {
		return fmt.Errorf("error while creating file : %w", err)
	}
	return nil
}

func (s *FileService) UpdateFile(ctx context.Context, file *entities.File) error {
	if err := s.db.Files.UpdateAndUpdateQuotaAndVersion(ctx, file); err != nil {
		return fmt.Errorf("error while updating file : %w", err)
	}
	return nil
}

func (s *FileService) GetFileById(ctx context.Context, id string) (*entities.File, error) {
	file, err := s.db.Files.GetFileById(ctx, id)
	if err != nil {
		return nil, err
	}
	downloadUrl, err := s.s3Service.Get(ctx, file.RemoteFileName)
	file.DownloadUrl = downloadUrl
	if err != nil {
		return nil, err
	}
	return file, nil
}
