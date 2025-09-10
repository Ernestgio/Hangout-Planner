package constants

const (

	// Environment constants
	ProductionEnv  = "PROD"
	DevEnv         = "DEV"
	DefaultAppName = "Hangout"
	DefaultAppPort = "9000"
	DefaultDBHost  = "localhost"
	DefaultDBPort  = "3306"
	DefaultDBUser  = "root"
	DefaultDBName  = "hangout"

	// Application constants
	LoggerFormat = "[${time_rfc3339}] ${method} ${host}${uri} ${status} ${latency_human}\n"

	//Status constants
	SuccessStatus    = "success"
	ErrorStatus      = "error"
	ProdErrorMessage = "An unexpected error occurred. Please try again later."

	// message constants
	UserCreatedSuccessfully = "User created successfully."

	// Error message constants
)
