package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// RunFullTest performs comprehensive MongoDB test
func RunFullTest() bool {
	log.Println("🧪 Starting comprehensive MongoDB test...")

	// Test 1: Connection test
	if !TestConnection() {
		return false
	}

	// Test 2: Insert test document
	collection := MongoDb.Collection("conversations")

	testConv := Conversation{
		AppID:       "test_app_" + fmt.Sprintf("%d", time.Now().Unix()),
		UserToken:   "test_user_123",
		SessionID:   "sess_abc",
		UserMessage: "Hello, is this working?",
		BotReply:    "Yes! Connection successful.",
		Metadata: map[string]interface{}{
			"test_type": "connection_test",
			"timestamp": time.Now().Unix(),
		},
		Timestamp: time.Now(),
	}

	insertResult, err := collection.InsertOne(context.Background(), testConv)
	if err != nil {
		log.Println("❌ Failed to insert test document:", err)
		return false
	}
	log.Printf("✅ Test document inserted with ID: %v\n", insertResult.InsertedID)

	// Test 3: Read back test document
	var result Conversation
	err = collection.FindOne(context.Background(), map[string]interface{}{"app_id": testConv.AppID}).Decode(&result)
	if err != nil {
		log.Println("❌ Failed to read test document:", err)
		return false
	}
	log.Printf("✅ Retrieved: User='%s', Message='%s'\n", result.UserToken, result.UserMessage)

	// Test 4: Cleanup test document
	_, err = collection.DeleteOne(context.Background(), map[string]interface{}{"app_id": testConv.AppID})
	if err != nil {
		log.Println("⚠️  Could not cleanup test document:", err)
	} else {
		log.Println("✅ Test document cleaned up")
	}

	log.Println("🎉 All tests passed successfully!")
	return true
}

// QuickTest performs a quick connection test only
func QuickTest() bool {
	log.Println("⚡ Performing quick connection test...")
	return TestConnection()
}
