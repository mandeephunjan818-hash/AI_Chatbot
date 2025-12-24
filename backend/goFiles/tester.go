// tester.go
package main

import (
	"context"
	"log"
	"time"
)

func tester() {
	InitMongoDB()
	defer func() {
		if err := MongoDb.Client().Disconnect(context.Background()); err != nil {
			log.Fatal("❌ Failed to disconnect from MongoDB:", err)
		}
	}()

	// Test: Insert a conversation
	collection := MongoDb.Collection("conversations")

	testConv := Conversation{
		AppID:       "test_app",
		UserToken:   "test_user_123",
		SessionID:   "sess_abc",
		UserMessage: "Hello, is this working?",
		BotReply:    "Yes! Connection successful.",
		Metadata: map[string]interface{}{
			"intent": "greeting",
			"mood":   "neutral",
		},
		Timestamp: time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), testConv)
	if err != nil {
		log.Fatal("❌ Failed to insert test document:", err)
	}

	log.Println("✅ Test conversation inserted successfully!")

	// Test: Read it back
	var result Conversation
	err = collection.FindOne(context.Background(), map[string]string{"app_id": "test_app"}).Decode(&result)
	if err != nil {
		log.Fatal("❌ Failed to read test document:", err)
	}

	log.Printf("✅ Retrieved: User='%s', Message='%s'\n", result.UserToken, result.UserMessage)
}
