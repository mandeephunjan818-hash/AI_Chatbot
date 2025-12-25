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

	MongoDb = client.Database(dbName)
	log.Println("✅ MongoDB client initialized. Database:", dbName)
}

// TestConnection tests the MongoDB connection
func TestConnection() bool {
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := MongoDb.Client().Ping(pingCtx, nil); err != nil {
		log.Println("❌ MongoDB ping failed:", err)
		return false
	}
	log.Println("✅ MongoDB connection test successful!")
	return true
}

// GetConnectionInfo returns connection details
func GetConnectionInfo() map[string]string {
	return map[string]string{
		"database": MongoDb.Name(),
		"uri":      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		"status":   "connected",
	}
}

// getEnv safely reads env var with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
