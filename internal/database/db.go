package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("❌ MONGO_URI environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}

	Client = client
	log.Println("✅ MongoDB connected successfully!")
}

func GetCollection(databaseName, collectionName string) *mongo.Collection {
	if Client == nil {
		log.Fatal("❌ MongoDB client is not initialized. Call Connect() first.")
	}
	return Client.Database(databaseName).Collection(collectionName)
}
