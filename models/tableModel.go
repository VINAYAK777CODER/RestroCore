package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
    ID               primitive.ObjectID `bson:"_id" json:"id"`
    Number_of_guests *int               `json:"number_of_guests" binding:"required" validate:"required,gte=1,lte=20" example:"4"`
    Table_number     *int               `json:"table_number" binding:"required" validate:"required,gte=1,lte=100" example:"12"`
    Created_at       time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
    Updated_at       time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
    Table_id         string             `bson:"table_id" json:"table_id" example:"TBL001"`
}


