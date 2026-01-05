package logmsg

// Application Lifecycle
const (
	AppExitSuccess         = "application exited successfully"
	AppCreateFailed        = "failed to create application"
	AppTerminatedWithError = "application terminated with error"
)

// Shutdown Messages
const (
	ShutdownSignalReceived  = "received shutdown signal"
	ShutdownInitiating      = "initiating graceful shutdown"
	ShutdownComplete        = "shutdown complete"
	ShutdownTimeoutExceeded = "graceful shutdown timeout exceeded, forcing stop"
)

// Database Messages
const (
	DBConnectionFailed = "database connection failed"
	DBCloseFailed      = "failed to close database connection"
)

// Network & gRPC Server
const (
	NetworkListenerFailed     = "failed to create network listener"
	GRPCServerListening       = "gRPC server listening"
	GRPCServerError           = "gRPC server error"
	ServerTerminatedWithError = "server terminated with error"
)

// Configuration
const (
	ConfigLoadFailed = "failed to load configuration"
)
