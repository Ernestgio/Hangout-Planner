package request

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func BindAndValidate[T any](c echo.Context) (*T, error) {
	var req T

	if err := c.Bind(&req); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return nil, err
	}

	return &req, nil
}
