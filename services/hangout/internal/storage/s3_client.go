package storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"io"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	client             *s3.Client
	awsConfig          aws.Config
	bucketName         string
	externalEndpoint   string
	presignedURLExpiry time.Duration
}

func NewS3Client(ctx context.Context, cfg *config.S3Config) (*S3Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)

	if err != nil {
		return nil, apperrors.ErrFailedCreateS3Client
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = cfg.UsePathStyle
	})

	return &S3Client{
		client:             client,
		awsConfig:          awsCfg,
		bucketName:         cfg.BucketName,
		externalEndpoint:   cfg.ExternalEndpoint,
		presignedURLExpiry: cfg.GetPresignedURLExpiry(),
	}, nil
}

func (s *S3Client) Upload(ctx context.Context, path string, reader io.Reader, contentType string) error {

	content, err := io.ReadAll(reader)
	if err != nil {
		return apperrors.ErrFileUploadFailed
	}

	contentMD5 := calculateMD5Checksum(content)

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(s.bucketName),
		Key:                  aws.String(path),
		Body:                 bytes.NewReader(content),
		ContentType:          aws.String(contentType),
		ContentMD5:           aws.String(contentMD5),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	})

	if err != nil {
		return apperrors.ErrFileUploadFailed
	}

	return nil
}

func calculateMD5Checksum(content []byte) string {
	hash := md5.Sum(content)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func (s *S3Client) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return apperrors.ErrFileDeleteFailed
	}

	return nil
}

func (s *S3Client) GeneratePresignedURL(ctx context.Context, path string) (string, error) {
	externalClient := s3.NewFromConfig(s.awsConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s.externalEndpoint)
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(externalClient)

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = s.presignedURLExpiry
	})

	if err != nil {
		return "", apperrors.ErrGetPresignedURLFailed
	}

	return req.URL, nil
}
