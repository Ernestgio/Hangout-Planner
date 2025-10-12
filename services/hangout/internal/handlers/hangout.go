package handlers

import (
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
	"github.com/labstack/echo/v4"
)

type HangoutHandler interface {
	CreateHangout(c echo.Context) error
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
