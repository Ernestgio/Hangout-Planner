package storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"io"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
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
		return nil, apperrors.ErrFailedLoadAWSConfig
	}

	otelaws.AppendMiddlewares(&awsCfg.APIOptions)

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
		return apperrors.ErrFileReadFailed
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

func (s *S3Client) GeneratePresignedDownloadURL(ctx context.Context, path string) (string, error) {
	presignClient := s.newPresignClient()

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	}, s.withPresignExpiry)
	if err != nil {
		return "", apperrors.ErrPresignedDownloadURLFailed
	}

	return req.URL, nil
}

func (s *S3Client) GeneratePresignedUploadURL(ctx context.Context, path string, contentType string) (string, error) {
	presignClient := s.newPresignClient()

	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(path),
		ContentType: aws.String(contentType),
	}, s.withPresignExpiry)
	if err != nil {
		return "", apperrors.ErrPresignedUploadURLFailed
	}

	return req.URL, nil
}

func (s *S3Client) newPresignClient() *s3.PresignClient {
	externalClient := s3.NewFromConfig(s.awsConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s.externalEndpoint)
		o.UsePathStyle = true
	})
	return s3.NewPresignClient(externalClient)
}

func (s *S3Client) withPresignExpiry(opts *s3.PresignOptions) {
	opts.Expires = s.presignedURLExpiry
}

func calculateMD5Checksum(content []byte) string {
	hash := md5.Sum(content)
	return base64.StdEncoding.EncodeToString(hash[:])
}
