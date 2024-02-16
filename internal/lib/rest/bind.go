package rest

import (
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

var _validator = validator.New()

func Bind[T any](c *fiber.Ctx) (T, error) {
	var t T

	if err := c.BodyParser(&t); err != nil {
		return t, fiber.NewError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := _validator.Struct(t); err != nil {
		return t, fiber.NewError(http.StatusUnprocessableEntity, err.Error())
	}

	return t, nil
}
