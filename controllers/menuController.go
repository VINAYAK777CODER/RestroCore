package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/VINAYAK777CODER/RestroCore/database"
	"github.com/VINAYAK777CODER/RestroCore/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)
var menuCollection  *mongo.Collection=database.OpenCollection(database.Client,"menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// 2. Slice to hold all menus
		var allMenus []models.Menu

		// 3. Find all documents (empty filter)
		cursor, err := menuCollection.Find(ctx, bson.M{})
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

	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c* gin.Context){

	}
}

func UpdateMenu() gin.HandlerFunc{
	return func(c* gin.Context){
		
	}
}