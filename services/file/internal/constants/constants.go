package constants

const (
	// Application Config - Default environment variable values constants
	ProductionEnv  = "PROD"
	DevEnv         = "DEV"
	DefaultAppName = "file"
	DefaultAppPort = "9001"

	// Database Config - Default environment variable values constants
	DefaultDBHost = "localhost"
	DefaultDBPort = "3306"
	DefaultDBUser = "root"
	DefaultDBName = "file"

	// DB Config - Default values constants
	DefaultDBCharset = "utf8mb4"
	DefaultDBNetwork = "tcp"

	// Logger constants
	LoggerNotInitializedWarning = "logger not initialized, using default configuration"

	// File Upload Constants
	MaxFileSize                  = 10 * 1024 * 1024 // 10MB in bytes
	DefaultPresignedURLExpiryMin = 15

	// Application Timeouts
	GracefulShutdownTimeout = 10 // seconds

	// S3 Config - Default values constants
	DefaultS3Endpoint         = "http://localhost:4566"
	DefaultS3ExternalEndpoint = "http://localhost:4566"
	DefaultS3Region           = "ap-southeast-1"
	DefaultS3Bucket           = "hangout-files"
)

// AllowedFileExtensions contains the permitted file extensions for uploads
var AllowedFileExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

// AllowedMimeTypes maps extensions to their allowed MIME types
var AllowedMimeTypes = map[string][]string{
	".jpg":  {"image/jpeg"},
	".jpeg": {"image/jpeg"},
	".png":  {"image/png"},
	".gif":  {"image/gif"},
	".webp": {"image/webp"},
}
