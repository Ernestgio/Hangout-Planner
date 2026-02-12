package middlewares

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/otel"
	"github.com/labstack/echo/v4"
)

func MetricsMiddleware(metrics *otel.MetricsRecorder) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()

			path := c.Path()
			method := req.Method

			finish := metrics.StartRequest(ctx, extractDomain(path), extractOperation(path))

			err := next(c)

			status := "success"
			if err != nil || c.Response().Status >= 400 {
				status = "error"
			}
			finish(status)

			return err
		}
	}
}

func extractDomain(path string) string {
	if len(path) < 2 {
		return "unknown"
	}

	if path[0] == '/' {
		path = path[1:]
	}

	for i, c := range path {
		if c == '/' {
			return path[:i]
		}
	}

	return path
}

func extractOperation(path string) string {
	if len(path) < 2 {
		return "unknown"
	}

	if path[0] == '/' {
		path = path[1:]
	}

	slashIndex := -1
	for i, c := range path {
		if c == '/' {
			slashIndex = i
			break
		}
	}

	if slashIndex == -1 || slashIndex >= len(path)-1 {
		return "unknown"
	}

	operation := path[slashIndex+1:]

	if len(operation) > 0 && operation[0] == ':' {
		return "get"
	}

	return operation
}
