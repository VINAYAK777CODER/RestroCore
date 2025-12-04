package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/VINAYAK777CODER/RestroCore/database"
	"github.com/VINAYAK777CODER/RestroCore/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")


var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		// context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		// slice to hold all the foods
		var allFoods []models.Food

		// -------------------- PAGINATION LOGIC STARTS --------------------
		// Query params: page and limit, if not provided default values are used
		pageStr := c.DefaultQuery("page", "1")  // page number
		limitStr := c.DefaultQuery("limit", "10") // how many items per page

		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)

		if page < 1 {
			page = 1 // page cannot be 0 or negative
		}
		if limit < 1 {
			limit = 10 // limit cannot be 0 or negative
		}

		skip := (page - 1) * limit // how many documents to skip
		// -------------------- PAGINATION LOGIC ENDS --------------------

		// find all the documents (empty filter)
		// added SetSkip and SetLimit for pagination
		cursor, err := foodCollection.Find(
			ctx,
			bson.M{},
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching foods"})
			return
		}
		defer cursor.Close(ctx) // always close the cursor

		// decode all the documents into the slice
		if err := cursor.All(ctx, &allFoods); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding food data"})
			return
		}

		// Successfully return paginated results
		c.JSON(http.StatusOK, allFoods)
	}
}


// Get Food
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

func round(num float64) int {

}

func toFixed(num float64, precision int) float64 {

}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {


	}
}
