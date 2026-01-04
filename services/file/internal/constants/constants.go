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

	// File constants
	DefaultPresignedURLExpiryMin = 15

	// S3 Config - Default values constants
	DefaultS3Endpoint         = "http://localhost:4566"
	DefaultS3ExternalEndpoint = "http://localhost:4566"
	DefaultS3Region           = "ap-southeast-1"
	DefaultS3Bucket           = "hangout-files"
)
