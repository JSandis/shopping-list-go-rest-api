package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jsandis/shopping-list-go-rest-api/controllers"
)

func ShoppingListRoute(app *fiber.App) {
	route := app.Group("/api/v1")

	route.Get("/shoppinglist", controllers.GetShoppingList)

	route.Get("/shoppinglist/item/:id", controllers.GetShoppingListItem)

	route.Post("/shoppinglist/item", controllers.AddShoppingListItem)

	route.Patch("/shoppinglist/item/:id", controllers.EditShoppingListItem)

	route.Delete("/shoppinglist/item/:id", controllers.DeleteShoppingListItem)
}
