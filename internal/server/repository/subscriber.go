package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriberRepository struct {
	col *mongo.Collection
}

func NewSubscriberRepository(col *mongo.Collection) *SubscriberRepository {
	return &SubscriberRepository{col: col}
}

func (r *SubscriberRepository) FindOne(ctx context.Context, filter bson.M) (*SubscriberState, error) {
	var state SubscriberState
	err := r.col.FindOne(ctx, filter).Decode(&state)
	return &state, err
}

func (r *SubscriberRepository) InsertOne(ctx context.Context, state *SubscriberState) error {
	_, err := r.col.InsertOne(ctx, state)
	return err
}

func (r *SubscriberRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

func (r *SubscriberRepository) DeleteOne(ctx context.Context, filter bson.M) error {
	_, err := r.col.DeleteOne(ctx, filter)
	return err
}
