package db

import (
	"context"
	"fmt"
	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
	"github.com/amankumarsingh77/cloudnest/internal/store/db"
)

type UserService struct {
	db *db.DbStore
}

func (s *UserService) CreateUser(ctx context.Context, user *entities.User) error {
	usernameExists, err := s.db.Users.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}
	if usernameExists != nil {
		return fmt.Errorf("username already used")
	}
	emailExists, err := s.db.Users.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	if emailExists != nil {
		return fmt.Errorf("email already used")
	}
	if err := s.db.Users.Create(ctx, user); err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	//TODO create user quota
	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *entities.User) error {
	if err := s.db.Users.Update(ctx, user); err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}
	return nil
}

func (s *UserService) GetUserById(ctx context.Context, userId string) (*entities.User, error) {
	user, err := s.db.Users.GetUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return user, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	user, err := s.db.Users.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := s.db.Users.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return user, nil
}
