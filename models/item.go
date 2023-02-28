package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShoppingListItem struct {
	Id        	primitive.ObjectID 	`json:"id" bson:"_id"`
	Name      	string             	`json:"name,omitempty" validate:"required"`
	Quantity  	int                	`json:"quantity,omitempty" validate:"required"`
	Status		bool				`json:"status,omitempty"`
	CreatedAt 	time.Time          	`json:"createdAt" bson:"createdAt"`
}
