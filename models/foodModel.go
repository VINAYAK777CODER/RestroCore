package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Food represents the structure of a food item in MongoDB.
type Food struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      *string            `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Price     *float64           `bson:"price" json:"price" validate:"required"`
	Food_image *string            `bson:"food_image" json:"food_image" validate:"required"`
	Created_at time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	Food_id    string             `bson:"food_id" json:"food_id"`
	Menu_id    *string            `bson:"menu_id" json:"menu_id" validate:"required"`
}
