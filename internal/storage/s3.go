package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/logger"
)

type Storage struct {
	client *s3.Client
}

var cfg = config.LoadConfig()

func NewStorage() StorageAdapter {
	awsCfg, _ := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithRegion(cfg.Region),
	)

	return &Storage{
		client: s3.NewFromConfig(awsCfg),
	}
}

func (s *Storage) SendObject(ctx context.Context, key, bucketname string) (ObjectUrl, error) {

	logger.Log.Info("starting s3 send object",
		"key", key,
	)

	presigner := s3.NewPresignClient(s.client)

	input := &s3.PutObjectInput{
		Bucket:      &bucketname,
		Key:         &key,
		ContentType: aws.String("image/jpeg"),
	}

	resp, err := presigner.PresignPutObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = 7 * 24 * time.Hour
	})

	if err != nil {
		logger.Log.Error("failed to generate presigned url", "err", err)
		return ObjectUrl{}, fmt.Errorf("presign put object: %w", err)
	}

	return ObjectUrl{
		UploadUrl: resp.URL,
		PhotoUrl:  fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketname, cfg.Region, key),
	}, nil
}
