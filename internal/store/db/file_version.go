package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
)

type VersionStore struct {
	db *sql.DB
}

func (s *VersionStore) CreateFileVersion(ctx context.Context, tx *sql.Tx, file *entities.File) error {
	query := `INSERT INTO file_versions (file_id, version_number, remote_file_name, size, checksum ) 
			  VALUES ($1, $2, $3, $4, $5)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(
		ctx,
		query,
		&file.ID,
		1,
		&file.RemoteFileName,
		&file.Size,
		&file.Checksum,
	)
	if err != nil {
		return fmt.Errorf("error creating file version: %w", err)
	}
	return nil
}

func (s *VersionStore) UpdateFileVersion(ctx context.Context, tx *sql.Tx, fileId string) error {
	query := `UPDATE file_versions SET version_number = version_number+1  Where file_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, fileId)
	if err != nil {
		return fmt.Errorf("error updating file version: %w", err)
	}
	return nil
}

func (s *VersionStore) GetFileVersion(ctx context.Context, fileId string) (*entities.FileVersion, error) {
	query := `SELECT id, file_id, version_number, remote_file_name, size, created_at FROM file_versions WHERE file_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	var fileVersion entities.FileVersion
	err := s.db.QueryRowContext(
		ctx,
		query,
		fileId,
	).Scan(
		&fileVersion.ID,
		&fileVersion.FileId,
		&fileVersion.VersionNum,
		&fileVersion.RemoteFileName,
		&fileVersion.Size,
		&fileVersion.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("error getting file version: %w", err)
	}
	return &fileVersion, err
}
