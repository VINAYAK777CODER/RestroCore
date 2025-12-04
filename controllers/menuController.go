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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var menuCollection  *mongo.Collection=database.OpenCollection(database.Client,"menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// 2. Slice to hold all menus
		var allMenus []models.Menu

		// -------------------- PAGINATION LOGIC START --------------------
		// Query parameters from client: page & limit
		pageStr := c.DefaultQuery("page", "1")   // which page user wants
		limitStr := c.DefaultQuery("limit", "10") // how many items per page

		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)

		// page and limit should never be zero or negative
		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}

		// skip = how many items to skip based on page
		skip := (page - 1) * limit
		// -------------------- PAGINATION LOGIC END ----------------------

		// 3. Find all documents (empty filter) + apply pagination
		cursor, err := menuCollection.Find(
			ctx, 
			bson.M{}, 
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while fetching menus",
			})
			return
		}
		defer cursor.Close(ctx) // 4. Always close cursor

		// 5. Decode all documents into slice
		if err := cursor.All(ctx, &allMenus); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while decoding menu data",
			})
			return
		}

		// 6. Send response
		c.JSON(http.StatusOK, allMenus)
	}
}




func GetMenu() gin.HandlerFunc{
	return func(c* gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c* gin.Context){
		var menu models.Menu
		ctx,cancel:=context.WithTimeout(context.Background(),time.Second*100)
		defer cancel()

		err:=c.BindJSON(&menu)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		validationErr:=validate.Struct(menu)
		if validationErr!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return 
		}

		// 5. Meta fields
		now := time.Now()
		menu.Created_at = now
		menu.Updated_at = now

		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		// 6. Insert into DB
		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			msg := fmt.Sprintf("menu is not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// 7. Success response
		c.JSON(http.StatusOK, result)




	}
}

// UpdateMenu updates an existing menu (partial update supported).
func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		if menuId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "menu_id is required in URL"})
			return
		}

		// Input struct for partial updates (all fields optional)
		type UpdateMenuInput struct {
			Name       *string    `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
			Category   *string    `json:"category,omitempty" validate:"omitempty,min=2,max=50"`
			Start_Date *time.Time `json:"start_date,omitempty"`
			End_Date   *time.Time `json:"end_date,omitempty"`
		}

		var input UpdateMenuInput

		// 1. Bind JSON
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2. If no fields provided, nothing to update
		if input.Name == nil &&
			input.Category == nil &&
			input.Start_Date == nil &&
			input.End_Date == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided to update"})
			return
		}

		// 3. Validate only provided fields
		if err := validate.Struct(input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 4. Business rule: if both dates are present, end_date must not be before start_date
		if input.Start_Date != nil && input.End_Date != nil {
			if input.End_Date.Before(*input.Start_Date) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "end_date cannot be before start_date",
				})
				return
			}
		}

		// 5. Build update object dynamically
		update := bson.D{}

		if input.Name != nil {
			update = append(update, bson.E{Key: "name", Value: *input.Name})
		}
		if input.Category != nil {
			update = append(update, bson.E{Key: "category", Value: *input.Category})
		}
		if input.Start_Date != nil {
			update = append(update, bson.E{Key: "start_date", Value: input.Start_Date})
		}
		if input.End_Date != nil {
			update = append(update, bson.E{Key: "end_date", Value: input.End_Date})
		}

		// Always update timestamp
		update = append(update, bson.E{Key: "updated_at", Value: time.Now()})

		filter := bson.M{"menu_id": menuId}

		// 6. Perform the update
		result, err := menuCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: update}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update menu",
			})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "menu not found",
			})
			return
		}

		// 7. Optionally return updated document
		var updatedMenu models.Menu
		if err := menuCollection.FindOne(ctx, filter).Decode(&updatedMenu); err != nil {
			// Update succeeded but fetch failed – still OK from a write perspective
			c.JSON(http.StatusOK, gin.H{
				"message": "menu updated successfully",
			})
			return
		}

		c.JSON(http.StatusOK, updatedMenu)
	}
}