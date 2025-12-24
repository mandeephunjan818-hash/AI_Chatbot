// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	InitMongoDB()
	defer func() {
		if err := MongoDb.Client().Disconnect(context.Background()); err != nil {
			log.Println("⚠️  Failed to disconnect from MongoDB:", err)
		}
	}()

	// Example HTTP handler (expand later)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("✅ API is running!\n"))
	})

	// Example: Log a dummy conversation (replace with real logic later)
	http.HandleFunc("/test-log", func(w http.ResponseWriter, r *http.Request) {
		collection := MongoDb.Collection("conversations")
		conv := Conversation{
			AppID:       "web_app",
			UserToken:   "user_xyz",
			SessionID:   "sess_" + time.Now().Format("20060102150405"),
			UserMessage: "Test from /test-log",
			BotReply:    "Logged successfully!",
			Metadata:    map[string]interface{}{"intent": "test", "mood": "neutral"},
			Timestamp:   time.Now(),
		}
		_, err := collection.InsertOne(context.Background(), conv)
		if err != nil {
			http.Error(w, "Failed to log", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("✅ Conversation logged to MongoDB!\n"))
	})

	log.Println("🚀 Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
