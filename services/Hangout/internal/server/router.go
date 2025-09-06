package server

import (
	"net/http"

	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/controllers"
	"github.com/labstack/echo/v4"
)

type Router struct {
	userController *controllers.UserController
}

func NewRouter(userController *controllers.UserController) *Router {
	return &Router{
		userController: userController,
	}
}

func (r *Router) RegisterEndpoints(server *echo.Echo) {
	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world! Hangout Planner API is running!")
	})

	server.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// User-related endpoints
	server.POST("/users", r.userController.CreateUser)
}
