package constants

const (
	// Environment constants
	ProductionEnv = "PROD"

	// Application constants
	LoggerFormat = "[${time_rfc3339}] ${method} ${host}${uri} ${status} ${latency_human}\n"

	//Status constants
	SuccessStatus = "success"
	ErrorStatus   = "error"

	// message constants
	UserCreatedSuccessfully = "User created successfully."
)
