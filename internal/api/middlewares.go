package api

import (
	"bytes"
	"encoding/json"
	"github.com/lcnssantos/rinha-de-backend/internal/lib/logging"
	"io"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := uuid.NewString()

		l := logging.Info(c.Request().Context())

		l = l.Str("path", c.Path()).
			Str("method", c.Request().Method).
			Interface("headers", c.Request().Header).
			Str("requestId", requestId)

		payload := map[string]any{}

		rawBody, err := io.ReadAll(c.Request().Body)

		c.Request().Body = io.NopCloser(bytes.NewBuffer(rawBody))

		if err == nil {
			err := json.Unmarshal(rawBody, &payload)

			if err == nil {
				l.Interface("payload", payload).Send()
			}
		}

		l.Msg("new http request")

		return next(c)
	}
}
