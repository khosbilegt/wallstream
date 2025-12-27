package service

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type FileService struct {
	uploadDir string
}

func NewFileService(uploadDir string) *FileService {
	return &FileService{uploadDir: uploadDir}
}

func (s *FileService) UploadFileStream(ctx context.Context, file multipart.File, filename string) (string, error) {
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return "", err
	}

	uniqueName := uuid.New().String() + filepath.Ext(filename)
	dstPath := filepath.Join(s.uploadDir, uniqueName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return uniqueName, nil
}
