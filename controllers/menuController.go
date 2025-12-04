package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VINAYAK777CODER/RestroCore/database"
	"github.com/VINAYAK777CODER/RestroCore/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func UpdateMenu() gin.HandlerFunc{
	return func(c* gin.Context){
		
	}
}