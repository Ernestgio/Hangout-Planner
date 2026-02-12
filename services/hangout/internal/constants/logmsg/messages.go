package logmsg

// Application Lifecycle
const (
	ConfigLoadFailed       = "Error loading config: %v"
	AppCreateFailed        = "Error creating app: %v"
	AppTerminatedWithError = "Fatal error running application: %v"
	AppExitSuccess         = "Application shutdown successfully"
)

// Shutdown Messages
const (
	AppShuttingDown              = "Received interrupt signal, shutting down..."
	DBConnectionCloseFailed      = "Error closing database connection: %v"
	FileServiceClientCloseFailed = "Error closing file service client: %v"
)

// gRPC Client
const (
	FileServiceClientInitialized = "File service client initialized: %s"
	FileServiceClientInitFailed  = "Failed to initialize file service client: %v"
)

// otel constants
const (
	OTELMeterProviderInitFailed = "Failed to initialize OTEL meter provider: %v"
	OTELRuntimeMetricsFailed    = "Failed to start OTEL runtime metrics: %v"
	OTELMetricsInitFailed       = "Failed to initialize OTEL metrics: %v"
	OTELShutdownFailed          = "Failed to shutdown OTEL provider: %v"
	OTELInitialized             = "OTEL initialized successfully"
)
