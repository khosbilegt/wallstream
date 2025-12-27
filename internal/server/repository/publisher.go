package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublisherRepository struct {
	col *mongo.Collection
}

func NewPublisherRepository(col *mongo.Collection) *PublisherRepository {
	return &PublisherRepository{col: col}
}

func (r *PublisherRepository) FindOne(ctx context.Context, filter bson.M) (*PublisherState, error) {
	var state PublisherState
	err := r.col.FindOne(ctx, filter).Decode(&state)
	return &state, err
}

func (r *PublisherRepository) InsertOne(ctx context.Context, state *PublisherState) error {
	_, err := r.col.InsertOne(ctx, state)
	return err
}

func (r *PublisherRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

func (r *PublisherRepository) DeleteOne(ctx context.Context, filter bson.M) error {
	_, err := r.col.DeleteOne(ctx, filter)
	return err
}
