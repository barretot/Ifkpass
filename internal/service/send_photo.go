package service

import (
	"context"
	"fmt"
	"time"

	"github.com/barretot/ifkpass/internal/storage"
)

type SendPhotoService struct {
	storage storage.StorageAdapter
}

func NewSendPhotoService(
	storage storage.StorageAdapter,
) *SendPhotoService {
	return &SendPhotoService{storage: storage}
}

func (s *SendPhotoService) SendPhoto(ctx context.Context, userID, bucketName string) (storage.ObjectUrl, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	key := fmt.Sprintf("users/%s/profile-photo.jpg", userID)

	output, err := s.storage.SendObject(ctx, key, bucketName)

	if err != nil {
		return storage.ObjectUrl{}, fmt.Errorf("storage: error send object: %w", err)
	}

	objectUrl := storage.ObjectUrl{
		UploadUrl: output.UploadUrl,
		PhotoUrl:  output.PhotoUrl,
	}

	return objectUrl, nil

}
