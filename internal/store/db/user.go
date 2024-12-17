package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *entities.User) error {
	query := `
	INSERT INTO users (username, password, email, name) VALUES ($1, $2, $3, $4) RETURNING id, status, email_verified, created_at, updated_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		&user.Username,
		&user.Password.Hash,
		&user.Email,
		&user.Name,
	).Scan(
		&user.ID,
		&user.Status,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}
	return nil
}

func (s *UserStore) Update(ctx context.Context, user *entities.User) error {
	query := `
	UPDATE users SET username = $1, email = $2, name = $3 WHERE id=$4 RETURNING status, email_verified, created_at, updated_at`

	res, err := s.db.ExecContext(
		ctx,
		query,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.ID,
	)
	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}
	if rows == 0 {
		return ErrorNotFound
	}
	return nil
}

func (s *UserStore) GetUserById(ctx context.Context, userId string) (*entities.User, error) {
	query := `
	SELECT id, name, email, username, status, email_verified, created_at, updated_at FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	var user entities.User
	err := s.db.QueryRowContext(
		ctx, query,
		&userId,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.Status,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, name, username, email, password, status, email_verified, last_login_at, created_at, updated_at FROM users WHERE email = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.Status,
		&user.EmailVerified,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get user by email: %w", err)
	}
	return &user, nil
}

func (s *UserStore) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, name, username, email, password, status, email_verified, last_login_at, created_at, updated_at FROM users WHERE username = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		username,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.Status,
		&user.EmailVerified,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get user by username: %s", username)
	}
	return &user, nil
}

func (s *UserStore) GetUserQuota(ctx context.Context, userId string) (*entities.UserQuota, error) {
	var quota entities.UserQuota
	query := `SELECT storage_used, storage_limit FROM user_quotas WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&quota.StorageUsed,
		&quota.StorageLimit,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("error getting quota: %w", err)
	}
	if quota.StorageUsed < quota.StorageLimit {
		quota.IsAllowedToUpload = true
	} else {
		quota.IsAllowedToUpload = false
	}
	return &quota, nil
}
