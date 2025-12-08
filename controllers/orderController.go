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
var orderCollection *mongo.Collection=database.OpenCollection(database.Client,"order")

func GetOrders() gin.HandlerFunc{
	return func(c* gin.Context){
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		//slice to hold all the orders
		var allOrders []models.Order

		cursor,err:=orderCollection.Find(ctx,bson.M{})
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occurred while fetching order"})
			return 
		}
		defer cursor.Close(ctx)

		if err:=cursor.All(ctx,&allOrders); err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while decoding the order data"})
		}

		c.JSON(http.StatusOK,allOrders)

		




	}
}

func GetOrder() gin.HandlerFunc{
	return func(c* gin.Context){
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		orderId:=c.Param("order_id")

		if orderId==""{
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
			return 
		}

		var order models.Order

		err:=orderCollection.FindOne(ctx,bson.M{"order_id":orderId}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusOK,order)


	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c* gin.Context){

	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c* gin.Context){
		
	}
}