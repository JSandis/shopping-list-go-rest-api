package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jsandis/shopping-list-go-rest-api/configs"
	"github.com/jsandis/shopping-list-go-rest-api/routes"
)

func main() {
	app := fiber.New()

	configs.ConnectDB()

	routes.ShoppingListRoute(app)
	routes.StatusNotFoundRoute(app)

	app.Listen(":" + configs.Port())
}
