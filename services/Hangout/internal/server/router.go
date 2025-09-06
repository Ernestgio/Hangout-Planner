package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
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

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world! Hangout Planner API is running!")
	})

	server.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// User routes
	usersRoute := server.Group("/users")
	usersRoute.POST("", r.dependencies.userController.CreateUser)
}
