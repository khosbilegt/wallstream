package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepository struct {
	col *mongo.Collection
}

func NewUsersRepository(col *mongo.Collection) *UsersRepository {
	return &UsersRepository{col: col}
}

func (r *UsersRepository) CreateUser(ctx context.Context, user *User) error {
	_, err := r.col.InsertOne(ctx, user)
	return err
}

func (r *UsersRepository) LoginUser(ctx context.Context, username, password string) (*User, error) {
	var user User
	err := r.col.FindOne(ctx, bson.M{"username": username, "password": password}).Decode(&user)
	return &user, err
}

func (r *UsersRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.col.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	return &user, err
}

func (r *UsersRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.col.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found → return nil, no error
			return nil, nil
		}
		// Other errors (connection, decode, etc.)
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetUserByAPIKey(ctx context.Context, apiKey string) (*User, error) {
	var user User
	err := r.col.FindOne(ctx, bson.M{"api_key": apiKey}).Decode(&user)
	return &user, err
}

func (r *UsersRepository) UpdateUserAPIKey(ctx context.Context, userID, apiKey string) error {
	_, err := r.col.UpdateOne(
		ctx,
		bson.M{"id": userID},
		bson.M{"$set": bson.M{"api_key": apiKey}},
	)
	return err
}
