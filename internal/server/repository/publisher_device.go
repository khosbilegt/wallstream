package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublisherDeviceRepository struct {
	col *mongo.Collection
}

func NewPublisherDeviceRepository(col *mongo.Collection) *PublisherDeviceRepository {
	return &PublisherDeviceRepository{col: col}
}

func (r *PublisherDeviceRepository) CreatePublisherDevice(ctx context.Context, publisherDevice *PublisherDevice) error {
	_, err := r.col.InsertOne(ctx, publisherDevice)
	return err
}

func (r *PublisherDeviceRepository) GetPublisherDevicesByUserID(ctx context.Context, userID string) ([]*PublisherDevice, error) {
	var publisherDevices []*PublisherDevice
	cursor, err := r.col.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var publisherDevice PublisherDevice
		err := cursor.Decode(&publisherDevice)
		if err != nil {
			return nil, err
		}
		publisherDevices = append(publisherDevices, &publisherDevice)
	}
	return publisherDevices, err
}

func (r *PublisherDeviceRepository) GetPublisherDeviceByDeviceID(ctx context.Context, deviceID string) (*PublisherDevice, error) {
	var publisherDevice PublisherDevice
	err := r.col.FindOne(ctx, bson.M{"device_id": deviceID}).Decode(&publisherDevice)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &publisherDevice, err
}

func (r *PublisherDeviceRepository) DeletePublisherDeviceByDeviceID(ctx context.Context, deviceID string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"device_id": deviceID})
	return err
}
