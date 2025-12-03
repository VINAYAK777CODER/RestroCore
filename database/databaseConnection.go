package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func DBinstance() *mongo.Client {

	mongoDb := "mongodb://localhost:27017"
	fmt.Println("Connecting to MongoDB...")

	// 1️⃣ Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2️⃣ Directly connect to MongoDB (recommended way)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDb))
	if err != nil {
		log.Fatal(err)
	}

	// 3️⃣ Ping MongoDB to confirm connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func ConnectDB() {
	Client = DBinstance()
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("restaurant").Collection(collectionName)
}
