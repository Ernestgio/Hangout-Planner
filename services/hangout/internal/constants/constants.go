package constants

const (

	// Application Config - Default environment variable values constants
	ProductionEnv  = "PROD"
	DevEnv         = "DEV"
	DefaultAppName = "hangout"
	DefaultAppPort = "9000"

	// Database Config - Default environment variable values constants
	DefaultDBHost = "localhost"
	DefaultDBPort = "3306"
	DefaultDBUser = "root"
	DefaultDBName = "hangout"

	// JWT Config - Default environment variable values constants
	DefaultJWTExpirationHours = 1

	// DB Config - Default values constants
	DefaultDBCharset = "utf8mb4"
	DefaultDBNetwork = "tcp"

	// Application constants
	LoggerFormat = "[${time_rfc3339}] ${method} ${host}${uri} ${status} ${latency_human}\n"

	// Application Timeouts
	GracefulShutdownTimeout = 10 // seconds

	// routes constants
	HealthCheckRoute = "/healthz"
	SwaggerRoute     = "/swagger/*"
	AuthRoutes       = "/auth"
	HangoutRoutes    = "/hangouts"
	ActivityRoutes   = "/activities"
	MemoryRoutes     = "/memories"

	//Status constants
	SuccessStatus = "success"
	ErrorStatus   = "error"
	HealthCheckOK = "OK"

	// message constants
	UserSignedUpSuccessfully = "User created successfully."
	UserSignedInSuccessfully = "User signed in successfully."

	HangoutCreatedSuccessfully    = "Hangout created successfully."
	HangoutUpdatedSuccessfully    = "Hangout updated successfully."
	HangoutRetrievedSuccessfully  = "Hangout retrieved successfully."
	HangoutDeletedSuccessfully    = "Hangout deleted successfully."
	HangoutsRetrievedSuccessfully = "Hangouts retrieved successfully."

	ActivityCreatedSuccessfully     = "Activity created successfully."
	ActivityUpdatedSuccessfully     = "Activity updated successfully."
	ActivityRetrievedSuccessfully   = "Activity retrieved successfully."
	ActivityDeletedSuccessfully     = "Activity deleted successfully."
	ActivitiesRetrievedSuccessfully = "Activities retrieved successfully."

	// Error message constants
	ProdErrorMessage = "An unexpected error occurred. Please try again later."
	ErrBadRequest    = "Bad request"

	// pagination
	DefaultLimit      = 20
	MaxLimit          = 100
	SortDirectionAsc  = "asc"
	SortDirectionDesc = "desc"
	SortByCreatedAt   = "created_at"
	SortByDate        = "date"

	// File upload constants
	MaxFileSize                  = 10 * 1024 * 1024
	MaxRequestSize               = 10 * 1024 * 1024
	DefaultPresignedURLExpiryMin = 15
	AllowedImageExtension        = ".jpg,.jpeg,.png,.gif,.webp"
	MaxFilePerUpload             = 10

	// Memory message constants
	MemoriesUploadedSuccessfully    = "Memories uploaded successfully."
	MemoryCreatedSuccessfully       = "Memory created successfully."
	MemoryFetchedSuccessfully       = "Memory fetched successfully."
	MemoryRetrievedSuccessfully     = "Memory retrieved successfully."
	MemoriesListedSuccessfully      = "Memories listed successfully."
	MemoriesRetrievedSuccessfully   = "Memories retrieved successfully."
	MemoryDeletedSuccessfully       = "Memory deleted successfully."
	UploadURLsGeneratedSuccessfully = "Upload URLs generated successfully."
	UploadConfirmedSuccessfully     = "Upload confirmed successfully."

	// S3 Config - Default environment variable values constants
	DefaultS3Endpoint         = "http://localhost:4566"
	DefaultS3ExternalEndpoint = "http://localhost:4566"
	DefaultS3Region           = "ap-southeast-1"
	DefaultS3Bucket           = "hangout-files"

	// grpc client default configs
	DefaultFileServiceURL = "file:9001"
	DefaultMTLSCertPath   = "/app/certs/mtls/hangout-client.crt"
	DefaultMTLSKeyPath    = "/app/certs/mtls/hangout-client.key"
	DefaultMTLSCAPath     = "/app/certs/mtls/ca.crt"
)
