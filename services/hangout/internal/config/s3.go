package config

import (
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
)

type S3Config struct {
	Endpoint              string
	Region                string
	AccessKeyID           string
	SecretAccessKey       string
	BucketName            string
	UsePathStyle          bool
	PresignedURLExpiryMin int
}

func NewS3Config() *S3Config {
	return &S3Config{
		Endpoint:              getEnv("S3_ENDPOINT", constants.DefaultS3Endpoint),
		Region:                getEnv("S3_REGION", constants.DefaultS3Region),
		AccessKeyID:           getEnv("S3_ACCESS_KEY_ID", "test"),
		SecretAccessKey:       getEnv("S3_SECRET_ACCESS_KEY", "test"),
		BucketName:            getEnv("S3_BUCKET_NAME", constants.DefaultS3Bucket),
		UsePathStyle:          true,
		PresignedURLExpiryMin: getEnvInt("S3_PRESIGNED_URL_EXPIRY_MINUTES", constants.DefaultPresignedURLExpiryMin),
	}
}

func (c *S3Config) GetPresignedURLExpiry() time.Duration {
	return time.Duration(c.PresignedURLExpiryMin) * time.Minute
}
