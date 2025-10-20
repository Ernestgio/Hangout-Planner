package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/auth"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/request"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/response"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/sanitizer"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type HangoutHandler interface {
	CreateHangout(c echo.Context) error
	UpdateHangout(c echo.Context) error
	GetHangoutByID(c echo.Context) error
	DeleteHangout(c echo.Context) error
}

type hangoutHandler struct {
	hangoutService  services.HangoutService
	responseBuilder *response.Builder
}

func NewHangoutHandler(hangoutService services.HangoutService, responseBuilder *response.Builder) HangoutHandler {
	return &hangoutHandler{
		hangoutService:  hangoutService,
		responseBuilder: responseBuilder,
	}
}

// @Summary      Create Hangout
// @Description  Creates a new hangout for the authenticated user.
// @Tags         Hangouts
// @Accept       json
// @Produce      json
// @Param        hangout body dto.CreateHangoutRequest true "Hangout creation data"
// @Success      201 {object} response.StandardResponse{data=dto.HangoutDetailResponse} "Hangout created successfully"
// @Failure      400 {object} response.StandardResponse "Invalid request payload"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /hangouts/ [post]
func (h *hangoutHandler) CreateHangout(c echo.Context) error {
	req, err := request.BindAndValidate[dto.CreateHangoutRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidPayload))
	}

	sanitizedTitle := sanitizer.SanitizeString(strings.TrimSpace(req.Title))
	sanitizedDescriptionHTML, err := sanitizer.SanitizeMarkdown(*req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(apperrors.ErrSanitizeDescription))
	}

	req.Title = sanitizedTitle
	req.Description = &sanitizedDescriptionHTML

	ctx := c.Request().Context()
	userToken := c.Get("userId").(*jwt.Token)
	claims := userToken.Claims.(*auth.TokenCustomClaims)
	userID := claims.UserID

	hangout, err := h.hangoutService.CreateHangout(ctx, userID, req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}

	return c.JSON(http.StatusCreated, h.responseBuilder.Success(constants.HangoutCreatedSuccessfully, hangout))
}

// @Summary      Update Hangout
// @Description  Updates an existing hangout for the authenticated user.
// @Tags         Hangouts
// @Accept       json
// @Produce      json
// @Param        hangout_id path string true "Hangout ID"
// @Param        hangout body dto.UpdateHangoutRequest true "Hangout update data"
// @Success      200 {object} response.StandardResponse{data=dto.HangoutDetailResponse} "Hangout updated successfully"
// @Failure      400 {object} response.StandardResponse "Invalid request payload"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      404 {object} response.StandardResponse "resource not found"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /hangouts/{hangout_id} [put]
func (h *hangoutHandler) UpdateHangout(c echo.Context) error {
	req, err := request.BindAndValidate[dto.UpdateHangoutRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidPayload))
	}

	sanitizedTitle := sanitizer.SanitizeString(strings.TrimSpace(req.Title))
	sanitizedDescriptionHTML, err := sanitizer.SanitizeMarkdown(*req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(apperrors.ErrSanitizeDescription))
	}

	req.Title = sanitizedTitle
	req.Description = &sanitizedDescriptionHTML

	ctx := c.Request().Context()
	userToken := c.Get("userId").(*jwt.Token)
	claims := userToken.Claims.(*auth.TokenCustomClaims)
	userID := claims.UserID
	hangoutId, err := uuid.Parse(c.Param("hangout_id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidHangoutID))
	}

	hangout, err := h.hangoutService.UpdateHangout(ctx, hangoutId, userID, req)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, h.responseBuilder.Error(apperrors.ErrNotFound))
		}
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}

	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.HangoutUpdatedSuccessfully, hangout))
}

// @Summary      Get Hangout by ID
// @Description  Retrieves a hangout by its ID for the authenticated user.
// @Tags         Hangouts
// @Accept       json
// @Produce      json
// @Param        hangout_id path string true "Hangout ID"
// @Success      200 {object} response.StandardResponse{data=dto.HangoutDetailResponse} "Hangout retrieved successfully"
// @Failure      400 {object} response.StandardResponse "Invalid Hangout ID"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      404 {object} response.StandardResponse "resource not found"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /hangouts/{hangout_id} [get]
func (h *hangoutHandler) GetHangoutByID(c echo.Context) error {
	hangoutId, err := uuid.Parse(c.Param("hangout_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidHangoutID))
	}
	ctx := c.Request().Context()
	userToken := c.Get("userId").(*jwt.Token)
	claims := userToken.Claims.(*auth.TokenCustomClaims)
	userID := claims.UserID
	hangout, err := h.hangoutService.GetHangoutByID(ctx, hangoutId, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, h.responseBuilder.Error(apperrors.ErrNotFound))
		}
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}
	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.HangoutRetrievedSuccessfully, hangout))
}

// @Summary      Delete Hangout
// @Description  Deletes a hangout by its ID for the authenticated user.
// @Tags         Hangouts
// @Accept       json
// @Produce      json
// @Param        hangout_id path string true "Hangout ID"
// @Success      200 {object} response.StandardResponse "Hangout deleted successfully"
// @Failure      400 {object} response.StandardResponse "Invalid Hangout ID"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      404 {object} response.StandardResponse "resource not found"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /hangouts/{hangout_id} [delete]
func (h *hangoutHandler) DeleteHangout(c echo.Context) error {
	hangoutId, err := uuid.Parse(c.Param("hangout_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidHangoutID))
	}
	ctx := c.Request().Context()
	userToken := c.Get("userId").(*jwt.Token)
	claims := userToken.Claims.(*auth.TokenCustomClaims)
	userID := claims.UserID
	err = h.hangoutService.DeleteHangout(ctx, hangoutId, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, h.responseBuilder.Error(apperrors.ErrNotFound))
		}
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}
	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.HangoutDeletedSuccessfully, nil))
}
