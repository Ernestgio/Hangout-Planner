package controllers

import (
	"net/http"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/mappings"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService     services.UserService
	responseBuilder *dto.StandardResponseBuilder
}

func NewUserController(userService services.UserService, responseBuilder *dto.StandardResponseBuilder) *UserController {
	return &UserController{
		userService:     userService,
		responseBuilder: responseBuilder,
	}
}

// @Summary      Create user
// @Description  Create a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.UserCreateRequest  true  "User data"
// @Success      201   {object}  dto.StandardResponseBuilder
// @Failure      400   {object}  dto.StandardResponseBuilder
// @Router       /users [post]
func (uc *UserController) CreateUser(c echo.Context) error {
	var req dto.UserCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, uc.responseBuilder.NewErrorResponse(apperrors.ErrInvalidPayload))
	}

	user, err := uc.userService.CreateUser(req)
	if err != nil {
		switch err {
		case apperrors.ErrUserAlreadyExists:
			return c.JSON(http.StatusConflict, uc.responseBuilder.NewErrorResponse(err))
		default:
			return c.JSON(http.StatusInternalServerError, uc.responseBuilder.NewErrorResponse(err))
		}
	}

	return c.JSON(http.StatusCreated, uc.responseBuilder.NewSuccessResponse(constants.UserCreatedSuccessfully, mappings.UserToResponseDTO(user)))
}
