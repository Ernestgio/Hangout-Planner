package logger

import (
	"fmt"
	"time"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LoggerFunc(c echo.Context, v middleware.RequestLoggerValues) error {
	fmt.Printf(
		constants.LoggerFormat,
		v.StartTime.Format(time.RFC3339),
		v.Status,
		v.Latency,
		v.Method,
		v.URI,
		v.RoutePath,
		v.Protocol,
		v.RemoteIP,
		v.ContentLength,
		v.ResponseSize,
		v.Error,
	)
	return nil
}
