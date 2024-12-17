package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
)

type FileStore struct {
	db *sql.DB
}

// TODO still not sure to use triggers or do stuff here using transactions
func (s *FileStore) create(ctx context.Context, tx *sql.Tx, file *entities.File) error {
	query := `INSERT INTO files (name, path, folder_id, remote_file_name, mime_type, size, checksum, url, created_by) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			  RETURNING created_at, updated_at, id`
	var folderId interface{}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		&file.Name,
		&file.Path,
		&folderId,
		&file.RemoteFileName,
		&file.MimeType,
		&file.Size,
		&file.Checksum,
		&file.Url,
		&file.CreatedBy,
	).Scan(
		&file.CreatedAt,
		&file.UpdatedAt,
		&file.ID,
	)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	return nil
}

func (s *FileStore) update(ctx context.Context, tx *sql.Tx, file *entities.File) error {
	query := `
    UPDATE files 
    SET 
        name = $1,
        folder_id = $2,
        remote_file_name = $3,
        size = $4,
        checksum = $5,
        url = $6, 
        path = $7,
        is_deleted = $8, 
        deleted_at = $9, 
        permanent_deletion_at = $10, 
        updated_at = $11
    WHERE id = $12`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	res, err := tx.ExecContext(
		ctx,
		query,
		&file.Name,
		&file.FolderId,
		&file.RemoteFileName,
		&file.Size,
		&file.Checksum,
		&file.Url,
		&file.Path,
		file.IsDeleted,
		file.DeletedAt,
		file.PermanentDeletionAt,
		file.UpdatedAt,
		file.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating file: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rows == 0 {
		return ErrorNotFound
	}
	return nil

}

// TODO implement Delete
func (s *FileStore) delete(ctx context.Context, id string) error {
	return nil
}

func (s *FileStore) GetFileById(ctx context.Context, id string) (*entities.File, error) {
	query := `SELECT id, name, folder_id, remote_file_name, mime_type, size, checksum, url, created_by, path, is_deleted, deleted_at, permanent_deletion_at, last_accessed_at, created_at, updated_at FROM files WHERE id = $1`
	var file entities.File
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&file.ID,
		&file.Name,
		&file.FolderId,
		&file.RemoteFileName,
		&file.MimeType,
		&file.Size,
		&file.Checksum,
		&file.Url,
		&file.CreatedBy,
		&file.Path,
		&file.IsDeleted,
		&file.DeletedAt,
		&file.PermanentDeletionAt,
		&file.LastAccessedAt,
		&file.CreatedAt,
		&file.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("error getting file: %w", err)
	}
	return &file, nil

}

func (s *FileStore) updateQuota(ctx context.Context, tx *sql.Tx, file *entities.File) error {
	query := `UPDATE user_quotas SET storage_used = storage_used + $1 Where user_id = $2 AND storage_used + $1 <= storage_limit`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, file.Size, file.CreatedBy)
	if err != nil {
		return fmt.Errorf("error updating user quota: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rows == 0 {
		return ErrQuotaExceeded
	}
	return nil
}

func (s *FileStore) createUserQuota(ctx context.Context, tx *sql.Tx, userQuota *entities.UserQuota) error {
	// TODO implement this
	query := `INSERT INTO user_quotas (user_id, storage_used, storage_limit) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := tx.ExecContext(
		ctx,
		query,
		userQuota.UserID,
		userQuota.StorageUsed,
		userQuota.StorageLimit,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileStore) addVersion(ctx context.Context, tx *sql.Tx, file *entities.File) error {
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
		return fmt.Errorf("error adding version: %w", err)
	}
	return nil
}

func (s *FileStore) updateVersion(ctx context.Context, tx *sql.Tx, fileId string) error {
	query := `UPDATE file_versions SET version_number = version_number+1  Where file_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, fileId)
	if err != nil {
		return fmt.Errorf("error updating file version: %w", err)
	}
	return nil
}

func (s *FileStore) CreateAndUpdateQuotaAndVersion(ctx context.Context, file *entities.File) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		if err := s.create(ctx, tx, file); err != nil {
			return err
		}
		if err := s.updateQuota(ctx, tx, file); err != nil {
			return err
		}
		if err := s.addVersion(ctx, tx, file); err != nil {
			return err
		}
		return nil
	})
}

func (s *FileStore) UpdateAndUpdateQuotaAndVersion(ctx context.Context, file *entities.File) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		if err := s.update(ctx, tx, file); err != nil {
			return err
		}
		if err := s.updateQuota(ctx, tx, file); err != nil {
			return err
		}
		if err := s.updateVersion(ctx, tx, file.ID); err != nil {
			return err
		}
		return nil
	})
}
