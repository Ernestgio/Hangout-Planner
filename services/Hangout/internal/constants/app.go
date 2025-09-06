package constants

const (
	ProductionEnv = "PROD"

	LoggerFormat = "[${time_rfc3339}] ${method} ${host}${uri} ${status} ${latency_human}\n"

	//Status constants
	SuccessStatus = "success"
	ErrorStatus   = "error"

	// message constants
	UserCreatedSuccessfully = "User created successfully."
)
