package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// StartServer starts the HTTP server
func StartServer() {
	// Set up HTTP routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/test-db", testDBHandler)
	http.HandleFunc("/db-info", dbInfoHandler)
	http.HandleFunc("/test-log", testLogHandler)
	http.HandleFunc("/full-test", fullTestHandler)

	log.Println("🚀 Starting HTTP server on :8080")
	log.Println("📋 Available endpoints:")
	log.Println("   GET  /health    - Basic health check")
	log.Println("   GET  /test-db   - Quick database connection test")
	log.Println("   GET  /db-info   - Database connection information")
	log.Println("   GET  /full-test - Run full database test")
	log.Println("   POST /test-log  - Log a test conversation")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service": "AI Chatbot API",
		"status":  "running",
		"version": "1.0.0",
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ API is running!\n"))
}

func testDBHandler(w http.ResponseWriter, r *http.Request) {
	if QuickTest() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("✅ Database connection is healthy!\n"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("❌ Database connection failed!\n"))
	}
}

func dbInfoHandler(w http.ResponseWriter, r *http.Request) {
	info := GetConnectionInfo()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func fullTestHandler(w http.ResponseWriter, r *http.Request) {
	result := RunFullTest()

	response := map[string]interface{}{
		"timestamp": time.Now().UTC(),
		"test_type": "full_database_test",
		"success":   result,
	}

	w.Header().Set("Content-Type", "application/json")
	if result {
		w.WriteHeader(http.StatusOK)
		response["message"] = "All tests passed successfully"
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		response["message"] = "Some tests failed"
	}

	json.NewEncoder(w).Encode(response)
}

func testLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed\n"))
		return
	}

	collection := MongoDb.Collection("conversations")
	conv := Conversation{
		AppID:       "web_app",
		UserToken:   "user_" + time.Now().Format("20060102150405"),
		SessionID:   "sess_" + time.Now().Format("20060102150405"),
		UserMessage: "Test from API endpoint",
		BotReply:    "Logged successfully via API!",
		Metadata: map[string]interface{}{
			"intent":    "test",
			"mood":      "neutral",
			"source":    "api_test",
			"timestamp": time.Now().Unix(),
		},
		Timestamp: time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), conv)
	if err != nil {
		http.Error(w, "Failed to log conversation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Conversation logged to MongoDB",
		"user_id": conv.UserToken,
	})
}
