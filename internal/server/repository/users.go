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
	return &user, err
}
