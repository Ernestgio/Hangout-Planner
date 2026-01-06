package apperrors

import "errors"

var ErrAppPortRequired = errors.New("APP_PORT is required")
var ErrDbPasswordRequired = errors.New("DB_PASSWORD required in production")

var ErrFailedLoadAWSConfig = errors.New("failed to load AWS config")
var ErrFailedCreateS3Client = errors.New("failed to create S3 client")
var ErrFileReadFailed = errors.New("failed to read file content")
var ErrFileUploadFailed = errors.New("file upload failed")
var ErrFileDeleteFailed = errors.New("file deletion failed")
var ErrPresignedDownloadURLFailed = errors.New("failed to generate presigned download URL")
var ErrPresignedUploadURLFailed = errors.New("failed to generate presigned upload URL")
