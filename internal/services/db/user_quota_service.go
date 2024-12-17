package db

import (
	"context"
	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
	"github.com/amankumarsingh77/cloudnest/internal/store/db"
)

type UserQuotaService struct {
	db *db.DbStore
}

func (s *UserQuotaService) GetUserQuota(ctx context.Context, userId string) (*entities.UserQuota, error) {
	userQuota, err := s.db.Users.GetUserQuota(ctx, userId)
	if err != nil {
		return nil, err
	}
	return userQuota, nil
}

func (s *UserQuotaService) CreateUserQuota(ctx context.Context, userQuota *entities.UserQuota) error {
	// TODO Not sure if it should be implemented
	return nil
}

func (s *UserQuotaService) UpdateUserQuota(ctx context.Context, userQuota *entities.UserQuota) error {
	// TODO implement this
	return nil
}
