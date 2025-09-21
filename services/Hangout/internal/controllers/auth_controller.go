package controllers

import (
	"net/http"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mappings"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService     services.AuthService
	responseBuilder *dto.StandardResponseBuilder
}

func NewAuthController(authService services.AuthService, responseBuilder *dto.StandardResponseBuilder) *AuthController {
	return &AuthController{
		authService:     authService,
		responseBuilder: responseBuilder,
	}
}

// @Summary      Sign up
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      dto.SignUpRequest  true  "User sign up data"
// @Success      201   {object}  dto.StandardResponse
// @Failure      400   {object}  dto.StandardResponse
// @Router       /auth/signup [post]
func (ac *AuthController) SignUp(c echo.Context) error {
	var req dto.SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ac.responseBuilder.NewErrorResponse(apperrors.ErrInvalidPayload))
	}
	user, err := ac.authService.SignUser(&req)
	if err != nil {
		switch err {
		case apperrors.ErrUserAlreadyExists:
			return c.JSON(http.StatusConflict, ac.responseBuilder.NewErrorResponse(err))
		default:
			return c.JSON(http.StatusInternalServerError, ac.responseBuilder.NewErrorResponse(err))
		}
	}

	return c.JSON(http.StatusCreated, ac.responseBuilder.NewSuccessResponse("User signed up successfully", mappings.UserToResponseDTO(user)))
}
