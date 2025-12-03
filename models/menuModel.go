package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name       string             `json:"name" validate:"required,min=2,max=50"`
    Category   string             `json:"category" validate:"required,min=2,max=50"`
    Start_Date *time.Time         `json:"start_date,omitempty" bson:"start_date,omitempty"`
    End_Date   *time.Time         `json:"end_date,omitempty" bson:"end_date,omitempty"`
    Created_at time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
    Updated_at time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
    Menu_id    string             `json:"menu_id" bson:"menu_id" validate:"required"`
}
