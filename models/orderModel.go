package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Order_Date time.Time          `json:"order_date" bson:"order_date" validate:"required"`
    Created_at time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
    Updated_at time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
    Order_id   string             `json:"order_id" bson:"order_id" validate:"required"`
    Table_id   *string            `json:"table_id" bson:"table_id" validate:"required"`
}
