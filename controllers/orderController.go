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

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//slice to hold all the orders
		var allOrders []models.Order

		cursor, err := orderCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching order"})
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &allOrders); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while decoding the order data"})
		}

		c.JSON(http.StatusOK, allOrders)

	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderId := c.Param("order_id")

		if orderId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
			return
		}

		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusOK, order)

	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a timeout context for MongoDB operations
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderId := c.Param("order_id")
		if orderId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is missing"})
			return
		}

		type UpdateOrderInput struct {
			Order_Date *time.Time `json:"order_date,omitempty"`
			Table_id   *string    `json:"table_id,omitempty"`
		}

		var input UpdateOrderInput

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		updateFields := bson.D{}

		if input.Order_Date != nil {
			updateFields = append(updateFields, bson.E{"order_date", input.Order_Date})
		}
		if input.Table_id != nil {
			updateFields = append(updateFields, bson.E{"table_id", input.Table_id})
		}

		if len(updateFields) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided to update"})
			return
		}

		//always update update_at

		updateFields = append(updateFields, bson.E{"updated_at", time.Now()})

		result, err := orderCollection.UpdateOne(
			ctx,
			bson.M{"order_id": orderId},
			bson.M{"$set": updateFields},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "order updated successfully"})

	}
}
