package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
    ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Invoice_id       string             `json:"invoice_id" bson:"invoice_id" validate:"required"`
    Order_id         string             `json:"order_id" bson:"order_id" validate:"required"`
    Payment_method   *string            `json:"payment_method,omitempty" bson:"payment_method,omitempty" validate:"omitempty,oneof=CARD CASH"`
    Payment_status   *string            `json:"payment_status,omitempty" bson:"payment_status,omitempty" validate:"required,oneof=PENDING PAID"`
    Payment_due_date time.Time          `json:"payment_due_date,omitempty" bson:"payment_due_date,omitempty"`
    Created_at       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
    Updated_at       time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
