package storage

import "context"

type ObjectUrl struct {
	UploadUrl string
	PhotoUrl  string
}

type StorageAdapter interface {
	SendObject(ctx context.Context, key, bucketname string) (ObjectUrl, error)
}
