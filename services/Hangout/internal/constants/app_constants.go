package constants

const (

	// Appliation Environment constants
	ProductionEnv  = "PROD"
	DevEnv         = "DEV"
	DefaultAppName = "Hangout"
	DefaultAppPort = "9000"
	DefaultDBHost  = "localhost"
	DefaultDBPort  = "3306"
	DefaultDBUser  = "root"
	DefaultDBName  = "hangout"

	// DB Environment constants
	DefaultDBCharset = "utf8mb4"

	// Application constants
	LoggerFormat = "[${time_rfc3339}] ${method} ${host}${uri} ${status} ${latency_human}\n"

	//Status constants
	SuccessStatus = "success"
	ErrorStatus   = "error"
	HealthCheckOK = "OK"

	// message constants
	ProdErrorMessage        = "An unexpected error occurred. Please try again later."
	UserCreatedSuccessfully = "User created successfully."

	// Error message constants
)
