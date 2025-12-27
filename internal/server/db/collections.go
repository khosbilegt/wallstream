package db

import "go.mongodb.org/mongo-driver/mongo"

type Collections struct {
	Users           *mongo.Collection
	PublisherState  *mongo.Collection
	SubscriberState *mongo.Collection
	RefreshTokens   *mongo.Collection
}

func NewCollections(db *mongo.Database) *Collections {
	return &Collections{
		Users:           db.Collection("users"),
		PublisherState:  db.Collection("publisher_states"),
		SubscriberState: db.Collection("subscriber_states"),
		RefreshTokens:   db.Collection("refresh_tokens"),
	}
}
