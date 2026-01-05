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
	AppShuttingDown         = "Received interrupt signal, shutting down..."
	DBConnectionCloseFailed = "Error closing database connection: %v"
)
