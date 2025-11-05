package handlers

import (
	"net/http"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/request"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/response"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ActivityHandler interface {
	CreateActivity(c echo.Context) error
	GetActivityByID(c echo.Context) error
	GetAllActivities(c echo.Context) error
	UpdateActivity(c echo.Context) error
	DeleteActivity(c echo.Context) error
}

type activityHandler struct {
	activityService services.ActivityService
	responseBuilder *response.Builder
}

func NewActivityHandler(activityService services.ActivityService, responseBuilder *response.Builder) ActivityHandler {
	return &activityHandler{
		activityService: activityService,
		responseBuilder: responseBuilder,
	}
}

// @Summary      Create Activity
// @Description  Creates a new activity for the authenticated user.
// @Tags         Activities
// @Accept       json
// @Produce      json
// @Param        activity body dto.CreateActivityRequest true "Activity creation data"
// @Success      201 {object} response.StandardResponse{data=dto.ActivityDetailResponse} "Activity created successfully"
// @Failure      400 {object} response.StandardResponse "Invalid request payload"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /activities/ [post]
func (h *activityHandler) CreateActivity(c echo.Context) error {
	req, err := request.BindAndValidate[dto.CreateActivityRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(err))
	}

	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()
	activity, err := h.activityService.CreateActivity(ctx, userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}
	return c.JSON(http.StatusCreated, h.responseBuilder.Success(constants.ActivityCreatedSuccessfully, activity))
}

// @Summary      Get Activity by ID
// @Description  Retrieves an activity by its ID for the authenticated user.
// @Tags         Activities
// @Produce      json
// @Param        activity_id path string true "Activity ID"
// @Success      200 {object} response.StandardResponse{data=dto.ActivityDetailResponse} "Activity fetched successfully"
// @Failure      400 {object} response.StandardResponse "Invalid activity ID"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /activities/{activity_id} [get]
func (h *activityHandler) GetActivityByID(c echo.Context) error {
	activityIDParam := c.Param("activity_id")
	activityID, err := uuid.Parse(activityIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(err))
	}
	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()
	activity, err := h.activityService.GetActivityByID(ctx, activityID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}
	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.ActivityRetrievedSuccessfully, activity))
}

// @Summary      Get All Activities
// @Description  Retrieves all activities for the authenticated user.
// @Tags         Activities
// @Produce      json
// @Success      200 {object} response.StandardResponse{data=[]dto.ActivityListItemResponse} "Activities fetched successfully"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /activities/ [get]
func (h *activityHandler) GetAllActivities(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()
	activities, err := h.activityService.GetAllActivities(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}
	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.ActivitiesRetrievedSuccessfully, activities))
}

// @Summary      Update Activity
// @Description  Updates an existing activity for the authenticated user.
// @Tags         Activities
// @Accept       json
// @Produce      json
// @Param        activity_id path string true "Activity ID"
// @Param        activity body dto.UpdateActivityRequest true "Activity update data"
// @Success      200 {object} response.StandardResponse{data=dto.ActivityDetailResponse} "Activity updated successfully"
// @Failure      400 {object} response.StandardResponse "Invalid request payload"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      404 {object} response.StandardResponse "resource not found"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /activities/{activity_id} [put]
func (h *activityHandler) UpdateActivity(c echo.Context) error {
	activityIDParam := c.Param("activity_id")
	activityID, err := uuid.Parse(activityIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(err))
	}
	req, err := request.BindAndValidate[dto.UpdateActivityRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(err))
	}
	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()
	activity, err := h.activityService.UpdateActivity(ctx, activityID, userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}
	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.ActivityUpdatedSuccessfully, activity))
}

// @Summary      Delete Activity
// @Description  Deletes an existing activity for the authenticated user.
// @Tags         Activities
// @Produce      json
// @Param        activity_id path string true "Activity ID"
// @Success      200 {object} response.StandardResponse "Activity deleted successfully"
// @Failure      400 {object} response.StandardResponse "Invalid activity ID"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      404 {object} response.StandardResponse "resource not found"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /activities/{activity_id} [delete]
func (h *activityHandler) DeleteActivity(c echo.Context) error {
	activityIDParam := c.Param("activity_id")
	activityID, err := uuid.Parse(activityIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(err))
	}
	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()
	err = h.activityService.DeleteActivity(ctx, activityID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}
	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.ActivityDeletedSuccessfully, nil))
}
