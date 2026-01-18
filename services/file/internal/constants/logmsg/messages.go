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

// Storage Messages
const (
	S3ConnectionFailed = "S3 client initialization failed"
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

// OpenTelemetry
const (
	OTELInitialized    = "OpenTelemetry initialized"
	OTELInitFailed     = "failed to initialize OpenTelemetry"
	OTELShutdownFailed = "failed to shutdown OpenTelemetry tracer provider"
)

// mTLS
const (
	MTLSInitialized = "mTLS enabled for gRPC server"
	MTLSInitFailed  = "failed to initialize mTLS"
)
