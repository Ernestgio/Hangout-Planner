package constants

const (

	// Application Config - Devault environment variable values constants
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

	// routes constants
	HealthCheckRoute = "/healthz"
	SwaggerRoute     = "/swagger/*"
	AuthRoutes       = "/auth"
	HangoutRoutes    = "/hangouts"
	ActivityRoutes   = "/activities"

	//Status constants
	SuccessStatus = "success"
	ErrorStatus   = "error"
	HealthCheckOK = "OK"

	// message constants
	UserSignedUpSuccessfully      = "User created successfully."
	UserSignedInSuccessfully      = "User signed in successfully."
	HangoutCreatedSuccessfully    = "Hangout created successfully."
	HangoutUpdatedSuccessfully    = "Hangout updated successfully."
	HangoutRetrievedSuccessfully  = "Hangout retrieved successfully."
	HangoutDeletedSuccessfully    = "Hangout deleted successfully."
	HangoutsRetrievedSuccessfully = "Hangouts retrieved successfully."

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
)
