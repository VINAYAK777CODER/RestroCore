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
	Order_id    string             `json:"order_id"`
	Order_items []models.OrderItem `json:"order_items"`
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

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// get order_item_id from URL
		orderItemId := c.Param("order_item_id")
		if orderItemId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_item_id is required"})
			return
		}

		var orderItem models.OrderItem

		// correct bson field name
		err := orderItemCollection.FindOne(
			ctx,
			bson.M{"order_item_id": orderItemId},
		).Decode(&orderItem)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "order item not found"})
			return
		}

		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var pack OrderItemPack

		// 1️⃣ Bind request JSON
		if err := c.BindJSON(&pack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
			return
		}

		// 2️⃣ Validate order_id
		if pack.Order_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
			return
		}

		// 3️⃣ Validate order_items
		if len(pack.Order_items) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order items cannot be empty"})
			return
		}

		var itemsToInsert []interface{}

		// 4️⃣ Process each order item
		for i := range pack.Order_items {

			item := &pack.Order_items[i]

			// validate basic fields (food_id, quantity)
			if err := validate.StructPartial(item, "Food_id", "Quantity"); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// 5️⃣ Fetch food price from DB
			var food struct {
				Price float64 `bson:"price"`
			}

			err := foodCollection.FindOne(
				ctx,
				bson.M{"food_id": *item.Food_id},
			).Decode(&food)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid food_id: " + *item.Food_id,
				})
				return
			}

			// 6️⃣ Backend assigns unit_price
			var num=toFixed(*item.Unit_price,2)
			item.Unit_price = &num

			// 7️⃣ Backend-controlled fields
			item.ID = primitive.NewObjectID()
			item.Order_item_id = item.ID.Hex()
			item.Order_id = pack.Order_id
			item.Created_at = time.Now()
			item.Updated_at = time.Now()

			// 8️⃣ Final validation
			if err := validate.Struct(item); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			itemsToInsert = append(itemsToInsert, *item)
		}

		// 9️⃣ Insert all items together
		_, err := orderItemCollection.InsertMany(ctx, itemsToInsert)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order items"})
			return
		}

		// 🔟 Response
		c.JSON(http.StatusCreated, gin.H{
			"order_id":    pack.Order_id,
			"order_items": pack.Order_items,
		})
	}
}


func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
