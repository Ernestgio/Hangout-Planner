package config

import (
	"os"
	"testing"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/stretchr/testify/require"
)

func TestNewS3Config_TableDriven(t *testing.T) {
	orig := map[string]*string{}
	keys := []string{"S3_ENDPOINT", "S3_EXTERNAL_ENDPOINT", "S3_REGION", "S3_ACCESS_KEY_ID", "S3_SECRET_ACCESS_KEY", "S3_BUCKET_NAME", "S3_PRESIGNED_URL_EXPIRY_MINUTES"}
	for _, k := range keys {
		if v, ok := os.LookupEnv(k); ok {
			vv := v
			orig[k] = &vv
		} else {
			orig[k] = nil
		}
	}
	defer func() {
		for k, v := range orig {
			if v == nil {
				_ = os.Unsetenv(k)
			} else {
				_ = os.Setenv(k, *v)
			}
		}
	}()

	tests := []struct {
		name          string
		env           map[string]string
		wantExpiryMin int
	}{
		{name: "defaults", env: map[string]string{}, wantExpiryMin: constants.DefaultPresignedURLExpiryMin},
		{name: "custom values", env: map[string]string{"S3_ENDPOINT": "http://s3", "S3_EXTERNAL_ENDPOINT": "http://ext", "S3_REGION": "us-east-1", "S3_ACCESS_KEY_ID": "ak", "S3_SECRET_ACCESS_KEY": "sk", "S3_BUCKET_NAME": "b", "S3_PRESIGNED_URL_EXPIRY_MINUTES": "30"}, wantExpiryMin: 30},
		{name: "invalid expiry", env: map[string]string{"S3_PRESIGNED_URL_EXPIRY_MINUTES": "bad"}, wantExpiryMin: constants.DefaultPresignedURLExpiryMin},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k := range orig {
				_ = os.Unsetenv(k)
			}
			for k, v := range tt.env {
				_ = os.Setenv(k, v)
			}
			cfg := NewS3Config()
			if v, ok := tt.env["S3_ENDPOINT"]; ok && v != "" {
				require.Equal(t, v, cfg.Endpoint)
			} else {
				require.Equal(t, constants.DefaultS3Endpoint, cfg.Endpoint)
			}
			if v, ok := tt.env["S3_EXTERNAL_ENDPOINT"]; ok && v != "" {
				require.Equal(t, v, cfg.ExternalEndpoint)
			} else {
				require.Equal(t, constants.DefaultS3ExternalEndpoint, cfg.ExternalEndpoint)
			}
			if v, ok := tt.env["S3_REGION"]; ok && v != "" {
				require.Equal(t, v, cfg.Region)
			} else {
				require.Equal(t, constants.DefaultS3Region, cfg.Region)
			}
			if v, ok := tt.env["S3_ACCESS_KEY_ID"]; ok && v != "" {
				require.Equal(t, v, cfg.AccessKeyID)
			} else {
				require.Equal(t, "test", cfg.AccessKeyID)
			}
			if v, ok := tt.env["S3_SECRET_ACCESS_KEY"]; ok && v != "" {
				require.Equal(t, v, cfg.SecretAccessKey)
			} else {
				require.Equal(t, "test", cfg.SecretAccessKey)
			}
			if v, ok := tt.env["S3_BUCKET_NAME"]; ok && v != "" {
				require.Equal(t, v, cfg.BucketName)
			} else {
				require.Equal(t, constants.DefaultS3Bucket, cfg.BucketName)
			}
			require.True(t, cfg.UsePathStyle)
			require.Equal(t, tt.wantExpiryMin, cfg.PresignedURLExpiryMin)
			require.Equal(t, time.Duration(cfg.PresignedURLExpiryMin)*time.Minute, cfg.GetPresignedURLExpiry())
		})
	}
}
