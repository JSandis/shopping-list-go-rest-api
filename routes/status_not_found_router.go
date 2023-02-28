package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func StatusNotFoundRoute(app *fiber.App) {
	app.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": http.StatusNotFound,
				"error":  true,
				"msg":    "Sorry, the endpoint was not found",
			})
		},
	)
}
