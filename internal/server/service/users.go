package service

import (
	"context"
	"errors"

	"github.io/khosbilegt/wallstream/internal/server/repository"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repo *repository.UsersRepository
}

func NewUsersService(repo *repository.UsersRepository) *UsersService {
	return &UsersService{repo: repo}
}

func (s *UsersService) CreateUser(ctx context.Context, user *repository.User, rawPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return s.repo.CreateUser(ctx, user)
}

// LoginUser verifies the password against the stored hash
func (s *UsersService) LoginUser(ctx context.Context, username, rawPassword string) (*repository.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user.Password == "" {
		return nil, errors.New("user has no password set")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPassword)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
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
