package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type IS3Repository interface {
	UploadFile(ctx context.Context, key string, data io.Reader, size int64, contentType string) error
	DownloadFile(ctx context.Context, key string) ([]byte, error)
	GetPresignedURL(ctx context.Context, key string, expires time.Duration) (string, error)
	DeleteFile(ctx context.Context, key string) error
	EnsureBucket(ctx context.Context) error
}

type S3Repository struct {
	client *minio.Client
	bucket string
}

func NewS3Repository(client *minio.Client, bucket string) *S3Repository {
	return &S3Repository{client: client, bucket: bucket}
}

func (r *S3Repository) EnsureBucket(ctx context.Context) error {
	exists, err := r.client.BucketExists(ctx, r.bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		err = r.client.MakeBucket(ctx, r.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		logrus.Infof("Bucket '%s' created successfully", r.bucket)
	}
	return nil
}

func (r *S3Repository) UploadFile(ctx context.Context, key string, data io.Reader, size int64, contentType string) error {
	opts := minio.PutObjectOptions{ContentType: contentType}
	_, err := r.client.PutObject(ctx, r.bucket, key, data, size, opts)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

func (r *S3Repository) DownloadFile(ctx context.Context, key string) ([]byte, error) {
	obj, err := r.client.GetObject(ctx, r.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer obj.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, obj); err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}
	return buf.Bytes(), nil
}

func (r *S3Repository) GetPresignedURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := r.client.PresignedGetObject(ctx, r.bucket, key, expires, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presignedURL.String(), nil
}

func (r *S3Repository) DeleteFile(ctx context.Context, key string) error {
	return r.client.RemoveObject(ctx, r.bucket, key, minio.RemoveObjectOptions{})
}
