package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.io/khosbilegt/wallstream/internal/server/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublisherService struct {
	publisherRepo          *repository.PublisherDeviceRepository
	publishedWallpaperRepo *repository.PublishedWallpaperRepository
}

func NewPublisherService(publisherRepo *repository.PublisherDeviceRepository, publishedWallpaperRepo *repository.PublishedWallpaperRepository) *PublisherService {
	return &PublisherService{publisherRepo: publisherRepo, publishedWallpaperRepo: publishedWallpaperRepo}
}

func (s *PublisherService) CreatePublisherDevice(ctx context.Context, userID, deviceID string) error {
	publisherDevice := &repository.PublisherDevice{
		ID:        uuid.New().String(),
		UserID:    userID,
		DeviceID:  deviceID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	return s.publisherRepo.CreatePublisherDevice(ctx, publisherDevice)
}

// Generate url to upload the wallpaper to the server
func (s *PublisherService) GenerateUploadURL(ctx context.Context, userID, deviceID string) (string, error) {
	// TODO: Make it rely on the user ID and device ID to generate a unique upload URL
	// Generate a random string for the upload URL
	randomString := uuid.New().String()
	// Generate the upload URL
	uploadURL := fmt.Sprintf("https://wallstream.io/upload/%s", randomString)
	return uploadURL, nil
}

func (s *PublisherService) GetPublisherDevicesByUserID(ctx context.Context, userID string) ([]*repository.PublisherDevice, error) {
	return s.publisherRepo.GetPublisherDevicesByUserID(ctx, userID)
}

func (s *PublisherService) GetPublisherDeviceByDeviceID(ctx context.Context, deviceID string) (*repository.PublisherDevice, error) {
	publisherDevice, err := s.publisherRepo.GetPublisherDeviceByDeviceID(ctx, deviceID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	if publisherDevice == nil {
		return nil, nil
	}
	return publisherDevice, nil
}

func (s *PublisherService) DeletePublisherDeviceByDeviceID(ctx context.Context, deviceID string) error {
	return s.publisherRepo.DeletePublisherDeviceByDeviceID(ctx, deviceID)
}

// Publish wallpaper given file path that was already uploaded to the server
func (s *PublisherService) PublishUploadedWallpaper(ctx context.Context, userID, deviceID, filePath string) error {
	return nil
}
