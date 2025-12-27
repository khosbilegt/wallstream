package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.io/khosbilegt/wallstream/internal/server/repository"
)

type UsersService struct {
	repo *repository.UsersRepository
}

func NewUsersService(repo *repository.UsersRepository) *UsersService {
	return &UsersService{repo: repo}
}

func (s *UsersService) CreateUser(ctx context.Context, username string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if user != nil {
		return "", errors.New("user already exists")
	}

	apiKey, err := GenerateAPIKey()
	if err != nil {
		return "", err
	}

	user = &repository.User{
		ID:        uuid.New().String(),
		Username:  username,
		APIKey:    apiKey,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}

func (s *UsersService) GetUserByID(ctx context.Context, id string) (*repository.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UsersService) GetUserByUsername(ctx context.Context, username string) (*repository.User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}

func (s *UsersService) GetUserByAPIKey(ctx context.Context, apiKey string) (*repository.User, error) {
	return s.repo.GetUserByAPIKey(ctx, apiKey)
}

func (s *UsersService) UpdateUserAPIKey(ctx context.Context, userID, apiKey string) error {
	return s.repo.UpdateUserAPIKey(ctx, userID, apiKey)
}
