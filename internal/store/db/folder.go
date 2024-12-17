package db

import (
	"context"
	"database/sql"
	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
)

type FolderStore struct {
	db *sql.DB
}

func (s *FolderStore) Create(context context.Context, folder *entities.Folder) error {
	return nil
}

func (s *FolderStore) Delete(context context.Context, id string) error {
	return nil
}

func (s *FolderStore) GetFilesInFolder(context context.Context, id string) (*[]entities.File, error) {
	return nil, nil

}
