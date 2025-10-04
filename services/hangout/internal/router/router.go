package router

import (
	"net/http"

	_ "github.com/Ernestgio/Hangout-Planner/services/hangout/api"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/handlers"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewRouter(e *echo.Echo, authhandler handlers.AuthHandler) {
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Auth routes
	authRoute := e.Group("/auth")
	authRoute.POST("/signup", authhandler.SignUp)
	authRoute.POST("/signin", authhandler.SignIn)

}
