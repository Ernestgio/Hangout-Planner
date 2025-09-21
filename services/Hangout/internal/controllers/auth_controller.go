package controllers

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
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
