package apperrors

import "errors"

// application errors
var ErrAppPortRequired = errors.New("APP_PORT is required")
var ErrDbPasswordRequired = errors.New("DB_PASSWORD required in production")

// http errors
var ErrInvalidPayload = errors.New("invalid payload")

// business errors

// auth & user
var ErrUserAlreadyExists = errors.New("user with that email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotFound = errors.New("user not found")
var ErrUnauthorized = errors.New("Unauthorized")

//pagination error
var ErrInvalidCursorPagination = errors.New("invalid cursor pagination")

// hangout
var ErrForbidden = errors.New("forbidden")
var ErrNotFound = errors.New("resource not found")
var ErrSanitizeDescription = errors.New("failed to sanitize description")
var ErrInvalidHangoutID = errors.New("invalid hangout ID")
var ErrInvalidPagination = errors.New("invalid pagination")
var ErrInvalidActivityIDs = errors.New("one or more activity IDs are invalid or not found")

var ErrInvalidActivityID = errors.New("invalid activity ID")

// file & memory errors
var ErrFileTooLarge = errors.New("file size exceeds maximum allowed")
var ErrInvalidFileType = errors.New("invalid file type")
var ErrFileNotFound = errors.New("file not found")
var ErrDuplicateFileName = errors.New("file with this name already exists in this hangout")
var ErrInvalidMemoryID = errors.New("invalid memory ID")
var ErrTooManyFiles = errors.New("too many files")
var ErrAllFilesUploadFailed = errors.New("all files upload failed")

var ErrMemoryNotFound = errors.New("memory not found")

// s3 errors
var ErrFailedCreateS3Client = errors.New("failed to create S3 client")
var ErrFileUploadFailed = errors.New("file upload failed")
var ErrFileDeleteFailed = errors.New("file deletion failed")
var ErrGetPresignedURLFailed = errors.New("failed to get presigned URL")
