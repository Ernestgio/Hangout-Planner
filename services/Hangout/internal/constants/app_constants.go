package constants

const (

	// Appliation Environment constants
	ProductionEnv  = "PROD"
	DevEnv         = "DEV"
	DefaultAppName = "hangout"
	DefaultAppPort = "9000"
	DefaultDBHost  = "localhost"
	DefaultDBPort  = "3306"
	DefaultDBUser  = "root"
	DefaultDBName  = "hangout"

	// DB Environment constants
	DefaultDBCharset = "utf8mb4"
	DefaultDBNetwork = "tcp"

	// Application constants
	LoggerFormat = "[${time_rfc3339}] ${method} ${host}${uri} ${status} ${latency_human}\n"

	//Status constants
	SuccessStatus = "success"
	ErrorStatus   = "error"
	HealthCheckOK = "OK"

	// message constants
	UserCreatedSuccessfully = "User created successfully."

	// Error message constants
	ProdErrorMessage = "An unexpected error occurred. Please try again later."

	// jwt claims
	JwtUserIDClaim = "user_id"
)
