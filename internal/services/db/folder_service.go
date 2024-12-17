package db

import (
	"context"
	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
	"github.com/amankumarsingh77/cloudnest/internal/store/db"
)

type FolderService struct {
	db *db.DbStore
}

func (s *FolderService) CreateFolder(ctx context.Context, folder *entities.Folder) error {
	return nil
}

func (s *FolderService) DeleteFolder(context context.Context, id string) error {
	return nil
}

func (s *FolderService) UpdateFolder(ctx context.Context, folder *entities.Folder) error {
	return nil
}

func (s *FolderService) GetFilesInFolder(context context.Context, id string) (*[]entities.File, error) {
	return nil, nil
}
