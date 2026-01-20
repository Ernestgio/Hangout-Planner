package config

import "github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"

type GRPCClientConfig struct {
	FileServiceURL string
	MTLSEnabled    bool
	CertFile       string
	KeyFile        string
	CAFile         string
}

func NewGRPCClientConfig() *GRPCClientConfig {
	return &GRPCClientConfig{
		FileServiceURL: getEnv("FILE_SERVICE_URL", constants.DefaultFileServiceURL),
		MTLSEnabled:    getEnv("GRPC_MTLS_ENABLED", "true") == "true",
		CertFile:       getEnv("GRPC_MTLS_CERT_FILE", constants.DefaultMTLSCertPath),
		KeyFile:        getEnv("GRPC_MTLS_KEY_FILE", constants.DefaultMTLSKeyPath),
		CAFile:         getEnv("GRPC_MTLS_CA_FILE", constants.DefaultMTLSCAPath),
	}
}
