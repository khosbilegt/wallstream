package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublishedWallpaperRepository struct {
	col *mongo.Collection
}

func NewPublishedWallpaperRepository(col *mongo.Collection) *PublishedWallpaperRepository {
	return &PublishedWallpaperRepository{col: col}
}

func (r *PublishedWallpaperRepository) CreatePublishedWallpaper(ctx context.Context, publishedWallpaper *PublishedWallpaper) error {
	_, err := r.col.InsertOne(ctx, publishedWallpaper)
	return err
}

func (r *PublishedWallpaperRepository) GetPublishedWallpapersByUserID(ctx context.Context, userID string) ([]*PublishedWallpaper, error) {
	var publishedWallpapers []*PublishedWallpaper
	cursor, err := r.col.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var publishedWallpaper PublishedWallpaper
		err := cursor.Decode(&publishedWallpaper)
		if err != nil {
			return nil, err
		}
		publishedWallpapers = append(publishedWallpapers, &publishedWallpaper)
	}
	return publishedWallpapers, err
}

func (r *PublishedWallpaperRepository) GetPublishedWallpapersByDeviceID(ctx context.Context, deviceID string) ([]*PublishedWallpaper, error) {
	var publishedWallpapers []*PublishedWallpaper
	cursor, err := r.col.Find(ctx, bson.M{"device_id": deviceID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var publishedWallpaper PublishedWallpaper
		err := cursor.Decode(&publishedWallpaper)
		if err != nil {
			return nil, err
		}
		publishedWallpapers = append(publishedWallpapers, &publishedWallpaper)
	}
	return publishedWallpapers, err
}

func (r *PublishedWallpaperRepository) GetPublishedWallpaperByHash(ctx context.Context, hash string) (*PublishedWallpaper, error) {
	var publishedWallpaper PublishedWallpaper
	err := r.col.FindOne(ctx, bson.M{"hash": hash}).Decode(&publishedWallpaper)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &publishedWallpaper, nil
}

func (r *PublishedWallpaperRepository) DeletePublishedWallpaperByHash(ctx context.Context, hash string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"hash": hash})
	return err
}

func (r *PublishedWallpaperRepository) DeletePublishedWallpaperByDeviceID(ctx context.Context, deviceID string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"device_id": deviceID})
	return err
}
