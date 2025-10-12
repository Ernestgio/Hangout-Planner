package middlewares

import (
	"net/http"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/auth"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/response"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWT(cfg *config.Config, responseBuilder *response.Builder) echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.TokenCustomClaims)
		},
		SigningKey: []byte(cfg.JwtConfig.JWTSecret),
		ContextKey: "userId",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, responseBuilder.Error(apperrors.ErrUnauthorized))
		},
	}
	return echojwt.WithConfig(config)
}
