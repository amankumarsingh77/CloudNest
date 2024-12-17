package s3

import "context"

type S3Store interface {
	GetPresignedURL(ctx context.Context, filename string, size int64) (string, error)
	Delete(ctx context.Context, filename string) error
	Get(ctx context.Context, filename string) (string, error)
}
