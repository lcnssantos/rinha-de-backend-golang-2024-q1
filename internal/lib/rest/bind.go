package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Bind[T any](c echo.Context) (T, error) {
	var t T

	if err := c.Bind(&t); err != nil {
		return t, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(t); err != nil {
		return t, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return t, nil
}
