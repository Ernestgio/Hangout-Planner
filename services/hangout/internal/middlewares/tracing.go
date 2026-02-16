package middlewares

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// TracingMiddleware returns an Echo middleware that instruments HTTP requests with OpenTelemetry tracing
func TracingMiddleware(serviceName string) echo.MiddlewareFunc {
	return otelecho.Middleware(serviceName)
}
