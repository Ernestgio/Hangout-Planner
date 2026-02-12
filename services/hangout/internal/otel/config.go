package otel

type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	Endpoint       string
	UseStdout      bool
}
