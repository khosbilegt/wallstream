package db

import "go.mongodb.org/mongo-driver/mongo"

type Collections struct {
	Users               *mongo.Collection
	PublisherDevices    *mongo.Collection
	PublishedWallpapers *mongo.Collection
}

func NewCollections(db *mongo.Database) *Collections {
	return &Collections{
		Users:               db.Collection("users"),
		PublisherDevices:    db.Collection("publisher_devices"),
		PublishedWallpapers: db.Collection("published_wallpapers"),
	}
}
