package mys3

import (
	e "LostAndFound/internal/common/errors"
	storage_config "LostAndFound/internal/config/storage_config"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FileRepository struct {
	client  *s3.S3
	bucket  string
	baseURL string
}

func (f FileRepository) GeneratePresignedPutURL(key string, contentType string, expires time.Duration) (string, error) {
	req, _ := f.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(f.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})

	urlStr, err := req.Presign(expires)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return urlStr, nil
}

func (f FileRepository) DeleteFile(ctx context.Context, key string) error {
	_, err := f.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		var aerr awserr.Error
		if errors.As(err, &aerr) && aerr.Code() == s3.ErrCodeNoSuchKey {
			return e.ErrFileNotFound
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (f FileRepository) FileExists(ctx context.Context, key string) (bool, error) {
	_, err := f.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

func (f FileRepository) GetBaseURL() string {
	return f.baseURL
}

func (f FileRepository) GetBucket() string {
	return f.bucket
}

func NewFileStorage(client *s3.S3, config storage_config.S3Config) *FileRepository {
	return &FileRepository{
		client:  client,
		bucket:  config.Bucket,
		baseURL: config.Endpoint + "/" + config.Bucket,
	}
}
