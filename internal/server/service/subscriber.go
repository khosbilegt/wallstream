package service

import (
	"context"

	"github.io/khosbilegt/wallstream/internal/server/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type SubscriberService struct {
	repo repository.SubscriberRepository
}

func NewSubscriberService(repo repository.SubscriberRepository) *SubscriberService {
	return &SubscriberService{repo: repo}
}

func (s *SubscriberService) GetSubscriberState(ctx context.Context, id string) (*repository.SubscriberState, error) {
	return s.repo.FindOne(ctx, bson.M{"id": id})
}

func (s *SubscriberService) CreateSubscriberState(ctx context.Context, state *repository.SubscriberState) error {
	return s.repo.InsertOne(ctx, state)
}
