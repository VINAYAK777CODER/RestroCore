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
		pageStr := c.DefaultQuery("page", "1")    // page number
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

		// 🔥 Check if food_id is missing
		if foodId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "food_id is required"})
			return
		}

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

// UpdateFood updates one or more fields of an existing food item.
// Supports partial updates: name, price, food_image, and menu_id.
func UpdateFood() gin.HandlerFunc {
    return func(c *gin.Context) {

        // Create a timeout context for MongoDB operations
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Get food_id from URL parameter
        foodId := c.Param("food_id")
        if foodId == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "food_id is missing"})
            return
        }

        // Struct for handling optional update fields (PATCH-like behavior)
        type UpdateFoodInput struct {
            Name        *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
            Price       *float64 `json:"price,omitempty" validate:"omitempty"`
            Food_image  *string  `json:"food_image,omitempty" validate:"omitempty"`
            Menu_id     *string  `json:"menu_id,omitempty" validate:"omitempty"`
        }

        var input UpdateFoodInput

        // Parse and bind JSON body into input struct
        if err := c.BindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Ensure at least one field was provided
        if input.Name == nil &&
            input.Price == nil &&
            input.Food_image == nil &&
            input.Menu_id == nil {

            c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided to update"})
            return
        }

        // Validate only the fields provided in the request
        if err := validate.Struct(input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // ------------------------------
        // Build MongoDB update document
        // ------------------------------
        update := bson.D{}

        if input.Name != nil {
            update = append(update, bson.E{Key: "name", Value: *input.Name})
        }

        if input.Price != nil {
            // Round price to 2 decimal places before saving
            rounded := toFixed(*input.Price, 2)
            update = append(update, bson.E{Key: "price", Value: rounded})
        }

        if input.Food_image != nil {
            update = append(update, bson.E{Key: "food_image", Value: *input.Food_image})
        }

        if input.Menu_id != nil {
            update = append(update, bson.E{Key: "menu_id", Value: *input.Menu_id})
        }

        // Always update 'updated_at' timestamp
        update = append(update, bson.E{Key: "updated_at", Value: time.Now()})

        // Filter to match the correct food item
        filter := bson.M{"food_id": foodId}

        // Perform database update using $set operator
        result, err := foodCollection.UpdateOne(
            ctx,
            filter,
            bson.D{{Key: "$set", Value: update}},
        )

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "failed to update food",
            })
            return
        }

        // If no document was matched, return not found
        if result.MatchedCount == 0 {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "food not found",
            })
            return
        }

        // Fetch updated food to return in response
        var updatedFood models.Food
        if err := foodCollection.FindOne(ctx, filter).Decode(&updatedFood); err != nil {
            // Fallback response if decoding fails
            c.JSON(http.StatusOK, gin.H{
                "message": "food updated successfully",
            })
            return
        }

        // Send fully updated food as response
        c.JSON(http.StatusOK, updatedFood)
    }
}



