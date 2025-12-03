package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Quantity     *string            `json:"Quantity" bson:"quantity,omitempty" validate:"required,oneof=S M L"`
    Unit_price   *float64           `json:"unit_price" bson:"unit_price,omitempty" validate:"required"`
    Created_at   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
    Updated_at   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
    Food_id      *string            `json:"food_id" bson:"food_id" validate:"required"`
    Order_item_id string            `json:"order_item_id" bson:"order_item_id" validate:"required"`
    Order_id     string             `json:"order_id" bson:"order_id" validate:"required"`
}
