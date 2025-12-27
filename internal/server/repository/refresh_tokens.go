package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshTokensRepository struct {
	col *mongo.Collection
}

func NewRefreshTokensRepository(col *mongo.Collection) *RefreshTokensRepository {
	return &RefreshTokensRepository{col: col}
}

func (r *RefreshTokensRepository) SaveRefreshToken(ctx context.Context, token *RefreshToken) error {
	_, err := r.col.InsertOne(ctx, token)
	return err
}

func (r *RefreshTokensRepository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	err := r.col.FindOne(ctx, bson.M{"token": token}).Decode(&refreshToken)
	return &refreshToken, err
}

func (r *RefreshTokensRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"token": token})
	return err
}

func (r *RefreshTokensRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	_, err := r.col.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}

// CleanupExpiredTokens removes expired refresh tokens
func (r *RefreshTokensRepository) CleanupExpiredTokens(ctx context.Context) error {
	now := time.Now().Unix()
	_, err := r.col.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lt": now}})
	return err
}
