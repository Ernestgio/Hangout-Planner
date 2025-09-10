package cmd

import (
	"net/http"

	_ "github.com/Ernestgio/Hangout-Planner/services/Hangout/docs"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/constants"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	dependencies *AppDependencies
}

func NewRouter(dependencies *AppDependencies) *Router {
	return &Router{
		dependencies: dependencies,
	}
}

func (r *Router) RegisterEndpoints(server *echo.Echo) {

	server.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, constants.HealthCheckOK)
	})

	// Swagger endpoint
	server.GET("/swagger/*", echoSwagger.WrapHandler)

	// User routes
	usersRoute := server.Group("/users")
	usersRoute.POST("", r.dependencies.userController.CreateUser)
}
