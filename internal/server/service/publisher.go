package service

import (
	"context"

	"github.io/khosbilegt/wallstream/internal/server/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type PublisherService struct {
	repo repository.PublisherRepository
}

func NewPublisherService(repo repository.PublisherRepository) *PublisherService {
	return &PublisherService{repo: repo}
}

func (s *PublisherService) GetPublisherState(ctx context.Context, id string) (*repository.PublisherState, error) {
	return s.repo.FindOne(ctx, bson.M{"id": id})
}

func (s *PublisherService) CreatePublisherState(ctx context.Context, state *repository.PublisherState) error {
	return s.repo.InsertOne(ctx, state)
}

func (s *PublisherService) UpdatePublisherState(ctx context.Context, id string, state *repository.PublisherState) error {
	return s.repo.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": state})
}

func (s *PublisherService) DeletePublisherState(ctx context.Context, id string) error {
	return s.repo.DeleteOne(ctx, bson.M{"id": id})
}
