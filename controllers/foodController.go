package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VINAYAK777CODER/RestroCore/database"
	"github.com/VINAYAK777CODER/RestroCore/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection=database.OpenCollection(database.Client,"food")
var validate=validator.New()
func GetFoods() gin.HandlerFunc{
	return func(c* gin.Context){

	}
}

//Get Food
func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
			return
		}

		c.JSON(http.StatusOK, food)
	}
}


func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		var food models.Food

		// 1. Parse JSON
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2. Validate struct
		if validationErr := validate.Struct(food); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// 3. Check menu exists
		if err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu); err != nil {
			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		// 4. Handle price safely
		if food.Price == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "price is required"})
			return
		}
		num := toFixed(*food.Price, 2)
		food.Price = &num

		// 5. Meta fields
		now := time.Now()
		food.Created_at = now
		food.Updated_at = now

		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()

		// 6. Insert into DB
		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := fmt.Sprintf("food item is not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// 7. Success response
		c.JSON(http.StatusOK, result)
	}
}


func round(num float64) int{

}

func toFixed(num float64,precision int) float64{

}


func UpdateFood() gin.HandlerFunc{
	return func(c* gin.Context){
		
	}
}

