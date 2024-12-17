package s3

import (
	"context"
	"log"
	"mime"
	"path/filepath"
	"time"

	"github.com/amankumarsingh77/cloudnest/internal/env"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	client          *s3.Client
	presignedClient *s3.PresignClient
	config          *env.StorageCfg
}

func NewS3Service(client *s3.Client, presignedClient *s3.PresignClient, config *env.StorageCfg) *S3Service {
	return &S3Service{
		client:          client,
		presignedClient: presignedClient,
		config:          config,
	}
}

func (s *S3Service) GetPresignedURL(ctx context.Context, filename string, filesize int64) (string, error) {
	extension := filepath.Ext(filename)
	contentType := mime.TypeByExtension(extension)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	res, err := s.presignedClient.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: &s.config.S3Bucket,
			Key:    &filename,
			//ContentType:   aws.String(contentType),
			//ContentLength: aws.Int64(filesize),
		},
		s3.WithPresignExpires(20*time.Minute),
	)
	if err != nil {
		return "", err
	}
	return res.URL, nil
}

func (s *S3Service) Delete(ctx context.Context, filename string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.config.S3Bucket,
		Key:    &filename,
	})
	return err
}

func (s *S3Service) Get(ctx context.Context, filename string) (string, error) {
	res, err := s.presignedClient.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: &s.config.S3Bucket,
			Key:    &filename,
		},
		func(options *s3.PresignOptions) {
			options.Expires = time.Hour
		},
	)
	if err != nil {
		return "", err
	}
	return res.URL, nil
}

func (s *S3Service) FileExist(ctx context.Context, filename string) bool {
	log.Println(s.config.S3Bucket, filename)
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.config.S3Bucket,
		Key:    &filename,
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
