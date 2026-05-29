package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"gritdemo/internal/config"
)

// Storage provides S3-compatible file storage operations.
type Storage struct {
	client *s3.Client
	bucket string
	cfg    config.StorageConfig
}

// New creates a new Storage instance using the given config.
// Works with AWS S3, MinIO, Cloudflare R2, and Backblaze B2.
func New(cfg config.StorageConfig) (*Storage, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if cfg.Endpoint != "" {
				return aws.Endpoint{
					URL:               cfg.Endpoint,
					HostnameImmutable: true,
					SigningRegion:     cfg.Region,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		},
	)

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true // Required for MinIO
	})

	// Verify bucket exists with a quick head request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.Bucket),
	})
	if err != nil {
		// Try to create the bucket
		_, createErr := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(cfg.Bucket),
		})
		if createErr != nil {
			return nil, fmt.Errorf("bucket %q not accessible and cannot be created: %w", cfg.Bucket, err)
		}
	}

	// Always ensure public-read policy so uploaded files are accessible via URL.
	// This is idempotent — safe to call on every startup.
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"AWS": ["*"]},
			"Action": ["s3:GetObject"],
			"Resource": ["arn:aws:s3:::%s/*"]
		}]
	}`, cfg.Bucket)

	_, _ = client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(cfg.Bucket),
		Policy: aws.String(policy),
	})

	return &Storage{
		client: client,
		bucket: cfg.Bucket,
		cfg:    cfg,
	}, nil
}

// Upload stores a file in the bucket at the given key.
func (s *Storage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("uploading %q: %w", key, err)
	}
	return nil
}

// Download retrieves a file from the bucket.
func (s *Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("downloading %q: %w", key, err)
	}
	return result.Body, nil
}

// Delete removes a file from the bucket.
func (s *Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("deleting %q: %w", key, err)
	}
	return nil
}

// GetURL returns the public URL for a stored file.
func (s *Storage) GetURL(key string) string {
	endpoint := strings.TrimRight(s.cfg.Endpoint, "/")
	// Encode each path segment individually to preserve forward slashes
	segments := strings.Split(key, "/")
	for i, seg := range segments {
		segments[i] = url.PathEscape(seg)
	}
	return fmt.Sprintf("%s/%s/%s", endpoint, s.bucket, strings.Join(segments, "/"))
}

// GetSignedURL returns a pre-signed URL valid for the given duration.
func (s *Storage) GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s.client)
	result, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(duration))
	if err != nil {
		return "", fmt.Errorf("generating signed URL for %q: %w", key, err)
	}
	return result.URL, nil
}

// PresignPutURL generates a pre-signed PUT URL for direct browser upload.
func (s *Storage) PresignPutURL(ctx context.Context, key, contentType string) (string, error) {
	presigner := s3.NewPresignClient(s.client)
	result, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(1*time.Hour))
	if err != nil {
		return "", fmt.Errorf("generating presigned PUT URL for %q: %w", key, err)
	}
	return result.URL, nil
}
