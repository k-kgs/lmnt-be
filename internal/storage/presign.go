// Package storage issues pre-signed upload URLs against Supabase Storage's
// S3-compatible API. The handler that uses this never receives or handles
// raw photo bytes — the caller uploads directly to the returned URL. Using
// the standard AWS SDK v2 presigner (not a Supabase-specific client) keeps
// this code portable to the documented Backblaze B2 fallback: only the
// endpoint/credentials would change, not this package.
package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"kayam-be/internal/config"
)

const presignExpiry = 5 * time.Minute

type Presigner struct {
	client *s3.PresignClient
	bucket string
}

func NewPresigner(ctx context.Context, cfg config.Config) (*Presigner, error) {
	if cfg.SupabaseS3Endpoint == "" || cfg.SupabaseStorageBucket == "" {
		return nil, fmt.Errorf("storage: SUPABASE_S3_ENDPOINT and SUPABASE_STORAGE_BUCKET must be set")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.SupabaseS3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.SupabaseS3AccessKeyID, cfg.SupabaseS3SecretAccessKey, "",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("storage: load aws config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.SupabaseS3Endpoint)
		o.UsePathStyle = true // required by Supabase's S3-compatible endpoint
	})

	return &Presigner{
		client: s3.NewPresignClient(s3Client),
		bucket: cfg.SupabaseStorageBucket,
	}, nil
}

// PresignPut returns a URL the caller can PUT the object to directly,
// valid for presignExpiry. The ContentType constraint means an upload with
// a mismatched Content-Type header will be rejected by Storage itself, not
// just by our own validation.
func (p *Presigner) PresignPut(ctx context.Context, objectKey, contentType string) (string, error) {
	req, err := p.client.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(presignExpiry))
	if err != nil {
		return "", fmt.Errorf("storage: presign put: %w", err)
	}
	return req.URL, nil
}
