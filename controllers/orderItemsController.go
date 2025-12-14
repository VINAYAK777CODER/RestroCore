package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/VINAYAK777CODER/RestroCore/database"
	"github.com/VINAYAK777CODER/RestroCore/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderItemPack struct {
	Order_id    string              `json:"order_id"`
	Order_items []models.OrderItem  `json:"order_items"`
}


var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")


func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := orderItemCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching order items"})
			return
		}
		defer cursor.Close(ctx)

		var orderItems []models.OrderItem
		if err := cursor.All(ctx, &orderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding order items"})
			return
		}

		c.JSON(http.StatusOK, orderItems)
	}
}

func ItemsByOrder(id string) (orderItems []primitive.M, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := orderItemCollection.Find(ctx, bson.M{"order_id": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &orderItems); err != nil {
		return nil, err
	}

	return orderItems, nil
}





// “Give me all order items where order_id matches”
func GetOrderItemByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

		orderId := c.Param("order_id")
		if orderId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
			return
		}

		orderItems, err := ItemsByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch order items"})
			return
		}

		if len(orderItems) == 0 {
			c.JSON(http.StatusOK, []primitive.M{})
			return
		}

		c.JSON(http.StatusOK, orderItems)
	}
}






func GetOrderItem() gin.HandlerFunc{
	return func(c* gin.Context){

	}
}

func CreateOrderItem() gin.HandlerFunc{
	return func(c* gin.Context){

	}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(c* gin.Context){
		
	}
}