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

func NewRouter(e *echo.Echo, cfg *config.Config, responseBuilder *response.Builder, authHandler handlers.AuthHandler, hangoutHandler handlers.HangoutHandler, activityHandler handlers.ActivityHandler, memoryHandler handlers.MemoryHandler, memoryHandlerV2 handlers.MemoryHandlerV2) {
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
	hangoutRoutes.Use(middlewares.UserContextMiddleware)
	hangoutRoutes.POST("/", hangoutHandler.CreateHangout)
	hangoutRoutes.PUT("/:hangout_id", hangoutHandler.UpdateHangout)
	hangoutRoutes.GET("/:hangout_id", hangoutHandler.GetHangoutByID)
	hangoutRoutes.DELETE("/:hangout_id", hangoutHandler.DeleteHangout)
	hangoutRoutes.POST("/list", hangoutHandler.GetHangoutsByUserID)

	// activity routes
	activityRoutes := e.Group(constants.ActivityRoutes)
	activityRoutes.Use(middlewares.JWT(cfg, responseBuilder))
	activityRoutes.Use(middlewares.UserContextMiddleware)
	activityRoutes.POST("/", activityHandler.CreateActivity)
	activityRoutes.PUT("/:activity_id", activityHandler.UpdateActivity)
	activityRoutes.GET("/:activity_id", activityHandler.GetActivityByID)
	activityRoutes.DELETE("/:activity_id", activityHandler.DeleteActivity)
	activityRoutes.GET("/", activityHandler.GetAllActivities)

	// memory routes (nested under hangouts for create/list)
	hangoutRoutes.POST("/:hangout_id/memories", memoryHandler.CreateMemories)
	hangoutRoutes.GET("/:hangout_id/memories", memoryHandler.ListMemories)
	hangoutRoutes.POST("/:hangout_id/memories/v2/upload-urls", memoryHandlerV2.GenerateUploadURLs)
	hangoutRoutes.POST("/:hangout_id/memories/v2/confirm-upload", memoryHandlerV2.ConfirmUpload)
	hangoutRoutes.GET("/:hangout_id/memories/v2", memoryHandlerV2.ListMemories)

	// memory routes (flat for single resource operations)
	memoryRoutes := e.Group(constants.MemoryRoutes)
	memoryRoutes.Use(middlewares.JWT(cfg, responseBuilder))
	memoryRoutes.Use(middlewares.UserContextMiddleware)
	memoryRoutes.GET("/:memory_id", memoryHandler.GetMemory)
	memoryRoutes.DELETE("/:memory_id", memoryHandler.DeleteMemory)

	// memory V2 routes (client-side upload) - single resource operations
	memoryRoutesV2 := e.Group(constants.MemoryRoutes + "/v2")
	memoryRoutesV2.Use(middlewares.JWT(cfg, responseBuilder))
	memoryRoutesV2.Use(middlewares.UserContextMiddleware)
	memoryRoutesV2.GET("/:memory_id", memoryHandlerV2.GetMemory)
	memoryRoutesV2.DELETE("/:memory_id", memoryHandlerV2.DeleteMemory)
}
