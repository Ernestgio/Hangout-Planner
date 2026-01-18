package config

type MTLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
	CAFile   string
}

func NewMTLSConfig() *MTLSConfig {
	enabled := getEnv("MTLS_ENABLED", "false") == "true"

	return &MTLSConfig{
		Enabled:  enabled,
		CertFile: getEnv("MTLS_CERT_FILE", "/app/certs/file-server.crt"),
		KeyFile:  getEnv("MTLS_KEY_FILE", "/app/certs/file-server.key"),
		CAFile:   getEnv("MTLS_CA_FILE", "/app/certs/ca.crt"),
	}
}
