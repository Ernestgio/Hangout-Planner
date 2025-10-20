package router

import (
	"net/http"

	_ "github.com/Ernestgio/Hangout-Planner/services/hangout/api"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/handlers"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/response"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/middlewares"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewRouter(e *echo.Echo, cfg *config.Config, responseBuilder *response.Builder, authHandler handlers.AuthHandler, hangoutHandler handlers.HangoutHandler) {
	e.GET(constants.HealthCheckRoute, func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET(constants.SwaggerRoute, echoSwagger.WrapHandler)

	// Auth routes
	authRoutes := e.Group(constants.AuthRoutes)
	authRoutes.POST("/signup", authHandler.SignUp)
	authRoutes.POST("/signin", authHandler.SignIn)

	// hangout routes
	hangoutRoutes := e.Group(constants.HangoutRoutes)
	hangoutRoutes.Use(middlewares.JWT(cfg, responseBuilder))
	hangoutRoutes.POST("/", hangoutHandler.CreateHangout)
	hangoutRoutes.PUT("/:hangout_id", hangoutHandler.UpdateHangout)
	hangoutRoutes.GET("/:hangout_id", hangoutHandler.GetHangoutByID)
	hangoutRoutes.DELETE("/:hangout_id", hangoutHandler.DeleteHangout)
}
