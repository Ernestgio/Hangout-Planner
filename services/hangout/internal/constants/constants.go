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
	MaxFilePerUpload = 10

	// Memory message constants
	MemoryRetrievedSuccessfully     = "Memory retrieved successfully."
	MemoriesRetrievedSuccessfully   = "Memories retrieved successfully."
	MemoryDeletedSuccessfully       = "Memory deleted successfully."
	UploadURLsGeneratedSuccessfully = "Upload URLs generated successfully."
	UploadConfirmedSuccessfully     = "Upload confirmed successfully."

	// grpc client default configs
	DefaultFileServiceURL = "file:9001"
	DefaultMTLSCertPath   = "/app/certs/mtls/hangout-client.crt"
	DefaultMTLSKeyPath    = "/app/certs/mtls/hangout-client.key"
	DefaultMTLSCAPath     = "/app/certs/mtls/ca.crt"

	// OTEL default configs
	DefaultOTELEndpoint       = "otelcollector:4317"
	DefaultOTELServiceVersion = "1.0.0"
	DefaultOTELUseStdout      = "false"
	DefaultTraceSampleRate    = 1.0
)
