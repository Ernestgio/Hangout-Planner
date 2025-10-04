package controllers

import (
	"net/http"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/request"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/response"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mappings"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthController interface {
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
}

type authController struct {
	authService     services.AuthService
	responseBuilder *response.Builder
}

func NewAuthController(authService services.AuthService, responseBuilder *response.Builder) AuthController {
	return &authController{
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
// @Success      201   {object}  response.StandardResponse
// @Failure      400   {object}  response.StandardResponse
// @Failure      409   {object}  response.StandardResponse
// @Failure      500   {object}  response.StandardResponse
// @Router       /auth/signup [post]
func (ac *authController) SignUp(c echo.Context) error {
	req, err := request.BindAndValidate[dto.SignUpRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ac.responseBuilder.Error(apperrors.ErrInvalidPayload))
	}
	user, err := ac.authService.SignUser(req)
	if err != nil {
		switch err {
		case apperrors.ErrUserAlreadyExists:
			return c.JSON(http.StatusConflict, ac.responseBuilder.Error(err))
		default:
			return c.JSON(http.StatusInternalServerError, ac.responseBuilder.Error(err))
		}
	}

	return c.JSON(http.StatusCreated, ac.responseBuilder.Success(constants.UserSignedUpSuccessfully, mappings.UserToResponseDTO(user)))
}

// @Summary      Sign in
// @Description  Authenticate a user and return a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      dto.SignInRequest  true  "User sign in data"
// @Success      200          {object}  response.StandardResponse
// @Failure      400          {object}  response.StandardResponse
// @Failure      401          {object}  response.StandardResponse
// @Failure      500          {object}  response.StandardResponse
// @Router       /auth/signin [post]
func (ac *authController) SignIn(c echo.Context) error {
	req, err := request.BindAndValidate[dto.SignInRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ac.responseBuilder.Error(apperrors.ErrInvalidPayload))
	}

	token, err := ac.authService.SignInUser(req)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidCredentials:
			return c.JSON(http.StatusUnauthorized, ac.responseBuilder.Error(err))
		default:
			return c.JSON(http.StatusInternalServerError, ac.responseBuilder.Error(err))
		}
	}
	return c.JSON(http.StatusOK, ac.responseBuilder.Success(constants.UserSignedInSuccessfully, token))
}
