package mys3

import (
	s "LostAndFound/internal/config/storage_config"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewS3Client(cfg s.S3Config) (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Region),
		Endpoint:         aws.String(cfg.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		DisableSSL:       aws.Bool(!cfg.UseSSL),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	client := s3.New(sess)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.Bucket),
	})
	if err != nil {
		_, createErr := client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(cfg.Bucket),
		})
		if createErr != nil {
			return nil, fmt.Errorf("bucket creation failed: %w", createErr)
		}
	}

	return client, nil
}
