package middlewares

import (
	"net/http"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func UserContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken, ok := c.Get("userId").(*jwt.Token)
		if !ok || userToken == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
		}

		claims, ok := userToken.Claims.(*auth.TokenCustomClaims)
		if !ok || claims == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
		}

		c.Set("user_id", claims.UserID)
		return next(c)
	}
}
