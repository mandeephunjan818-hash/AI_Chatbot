// connector.go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDb *mongo.Database

// InitMongoDB initializes MongoDB connection using .env
func InitMongoDB() {
	// Load .env file (optional)
	_ = godotenv.Load()

	uri := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getEnv("MONGODB_DB", "ai_chatbot")

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Failed to create MongoDB client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}

	// Verify connection
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingCancel()
	if err = client.Ping(pingCtx, nil); err != nil {
		log.Fatal("❌ MongoDB ping failed:", err)
	}

	MongoDb = client.Database(dbName)
	log.Println("✅ MongoDB connected to database:", dbName)
}

// getEnv safely reads env var with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
