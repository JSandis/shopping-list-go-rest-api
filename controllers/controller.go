package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jsandis/shopping-list-go-rest-api/configs"
	"github.com/jsandis/shopping-list-go-rest-api/models"
	"github.com/jsandis/shopping-list-go-rest-api/responses"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var shoppingListCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), configs.EnvMongoDBCollectionName())
var validate = validator.New()

func GetShoppingList(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var shoppingList []models.ShoppingListItem
	defer cancel()

	results, err := shoppingListCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.ShoppingListResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var shoppingListItem models.ShoppingListItem
		if err = results.Decode(&shoppingListItem); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				responses.ShoppingListResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
		shoppingList = append(shoppingList, shoppingListItem)
	}

	return c.Status(http.StatusOK).JSON(
		responses.ShoppingListResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": shoppingList}})
}

func GetShoppingListItem(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	shoppingListItemId := c.Params("id")
	var shoppingListItem models.ShoppingListItem
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(shoppingListItemId)

	err := shoppingListCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&shoppingListItem)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			responses.ShoppingListResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.ShoppingListResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": shoppingListItem}})
}

func AddShoppingListItem(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var shoppingListItem models.ShoppingListItem
	defer cancel()

	if err := c.BodyParser(&shoppingListItem); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.ShoppingListResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&shoppingListItem); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.ShoppingListResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newShoppingListItem := models.ShoppingListItem{
		Id:        primitive.NewObjectID(),
		Name:      shoppingListItem.Name,
		Quantity:  shoppingListItem.Quantity,
		Status:    false,
		CreatedAt: primitive.NewObjectID().Timestamp()}

	result, err := shoppingListCollection.InsertOne(ctx, newShoppingListItem)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.ShoppingListResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(
		responses.ShoppingListResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func EditShoppingListItem(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	shoppingListItemId := c.Params("id")
	var shoppingListItem models.ShoppingListItem
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(shoppingListItemId)

	if err := c.BodyParser(&shoppingListItem); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.ShoppingListResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&shoppingListItem); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.ShoppingListResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"name": shoppingListItem.Name, "quantity": shoppingListItem.Quantity, "status": shoppingListItem.Status}
	result, err := shoppingListCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			responses.ShoppingListResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	var updatedShoppingListItem models.ShoppingListItem
	if result.MatchedCount == 1 {
		err := shoppingListCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedShoppingListItem)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				responses.ShoppingListResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(
		responses.ShoppingListResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedShoppingListItem}})
}

func DeleteShoppingListItem(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	shoppingListItemId := c.Params("id")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(shoppingListItemId)
	result, err := shoppingListCollection.DeleteOne(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ShoppingListResponse{Status: http.StatusBadRequest, Message: "errror", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.ShoppingListResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Shopping list item with specified ID not found"}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.ShoppingListResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Shopping list item successfully deleted"}})
}
