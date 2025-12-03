package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/VINAYAK777CODER/RestroCore/database"
	"github.com/VINAYAK777CODER/RestroCore/middlewares"
	"github.com/VINAYAK777CODER/RestroCore/routes"

	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection

func main() {
	// 1️⃣ Connect to MongoDB FIRST
	database.ConnectDB()

	// 2️⃣ Now open the food collection AFTER connection is ready
	foodCollection = database.OpenCollection(database.Client, "food")

	// 3️⃣ Setup port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// 4️⃣ Create router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.Authentication())

	// 5️⃣ Register routes
	routes.UserRoutes(router)
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	// 6️⃣ Start server
	router.Run(":" + port)
}
