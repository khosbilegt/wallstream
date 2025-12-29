package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.io/khosbilegt/wallstream/internal/core/assets"
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

// Generate url to upload the wallpaper to the server
func (s *PublisherService) GenerateUploadURL(ctx context.Context, userID, deviceID string) (string, error) {
	// TODO: Make it rely on the user ID and device ID to generate a unique upload URL
	// Generate a random string for the upload URL
	randomString := uuid.New().String()
	// Generate the upload URL
	uploadURL := fmt.Sprintf("https://wallstream.io/upload/%s", randomString)
	return uploadURL, nil
}

// TODO: Cleanup previous files
// Publish wallpaper given file path that was already uploaded to the server
func (s *PublisherService) PublishUploadedWallpaper(ctx context.Context, userID, deviceID, filename string) error {
	publisherDevice, err := s.publisherRepo.GetPublisherDeviceByDeviceID(ctx, deviceID)
	if err != nil {
		return err
	}
	if publisherDevice == nil {
		return fmt.Errorf("no publisher device found for device %s", deviceID)
	}
	if publisherDevice.UserID != userID {
		return fmt.Errorf("publisher device not found for user %s", userID)
	}
	// Check hash of the file exists in the database
	filePath := "uploads/" + filename
	hash, err := assets.HashFile(filePath)

	if err != nil {
		return err
	}

	previousPublishedWallpapers, err := s.publishedWallpaperRepo.GetPublishedWallpapersByDeviceID(ctx, hash)
	if err != nil {
		return err
	}
	for _, previousPublishedWallpaper := range previousPublishedWallpapers {
		if previousPublishedWallpaper.Hash == hash {
			return fmt.Errorf("published wallpaper already exists for hash %s", hash)
		}
	}
	publishedWallpaper := &repository.PublishedWallpaper{
		ID:        uuid.New().String(),
		UserID:    userID,
		DeviceID:  deviceID,
		Hash:      hash,
		URL:       filePath,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	s.publishedWallpaperRepo.CreatePublishedWallpaper(ctx, publishedWallpaper)
	return nil
}

func (s *PublisherService) GetPublishedWallpapersByUserID(ctx context.Context, userID string) ([]*repository.PublishedWallpaper, error) {
	return s.publishedWallpaperRepo.GetPublishedWallpapersByUserID(ctx, userID)
}

func (s *PublisherService) GetPublishedWallpapersByDeviceID(ctx context.Context, userID string, deviceID string) ([]*repository.PublishedWallpaper, error) {
	publishedWallpapers, err := s.publishedWallpaperRepo.GetPublishedWallpapersByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	for _, publishedWallpaper := range publishedWallpapers {
		if publishedWallpaper.UserID != userID {
			return nil, fmt.Errorf("published wallpaper not found for user %s", userID)
		}
	}
	return publishedWallpapers, nil
}

func (s *PublisherService) DeletePublishedWallpaperByHash(ctx context.Context, userID, hash string) error {
	return s.publishedWallpaperRepo.DeletePublishedWallpaperByHash(ctx, hash)
}
