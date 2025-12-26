// main_tester.go - Main test runner
package mongoDB_tester

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"backend/gofiles/mongoDB"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Command line flags
	mongoURI     = flag.String("uri", "mongodb://localhost:27017", "MongoDB connection URI")
	databaseName = flag.String("db", "ai_chatbot", "Database name")
	testType     = flag.String("test", "all", "Test type: connection, schema, crud, all")
	createSchema = flag.Bool("create-schema", false, "Create database schema (collections and indexes)")
	cleanup      = flag.Bool("cleanup", false, "Clean up test data after tests")
	verbose      = flag.Bool("verbose", false, "Enable verbose logging")
	timeout      = flag.Duration("timeout", 30*time.Second, "Connection timeout")
)

func teater() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if *verbose {
		log.Println("🔍 Verbose mode enabled")
	}

	// Create MongoDB configuration
	config := &mongoDB.Config{
		URI:      *mongoURI,
		Database: *databaseName,
		Options: options.Client().
			SetConnectTimeout(*timeout).
			SetServerSelectionTimeout(*timeout).
			SetMaxPoolSize(10),
	}

	// Get connector instance
	connector := mongoDB.GetInstance(config)

	// Connect to MongoDB
	log.Printf("🔗 Connecting to MongoDB at: %s", *mongoURI)
	log.Printf("📁 Database: %s", *databaseName)

	err := connector.Connect()
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer connector.Disconnect()

	log.Println("✅ Connected to MongoDB successfully")

	// Create schema tester
	schemaTester, err := NewSchemaTester(connector)
	if err != nil {
		log.Fatalf("❌ Failed to create schema tester: %v", err)
	}

	// Run tests based on command
	switch *testType {
	case "connection":
		log.Println("🧪 Running Connection Tests Only")
		err = runConnectionTests(schemaTester)

	case "schema":
		log.Println("📁 Running Schema Tests Only")
		err = runSchemaTests(schemaTester)

	case "crud":
		log.Println("🔧 Running CRUD Tests Only")
		err = runCRUDTests(schemaTester)

	case "all":
		log.Println("🚀 Running All Tests")
		err = schemaTester.RunAllTests()

	default:
		log.Fatalf("❌ Unknown test type: %s. Use: connection, schema, crud, all", *testType)
	}

	if err != nil {
		log.Fatalf("❌ Tests failed: %v", err)
	}

	if *cleanup {
		log.Println("🧹 Cleaning up test data...")
		cleanupTestData(schemaTester)
	}

	log.Println("🎉 All operations completed successfully!")
}

func runConnectionTests(tester *SchemaTester) error {
	log.Println("🔗 Testing MongoDB Connection...")

	// Test basic connectivity
	err := tester.TestConnection()
	if err != nil {
		return fmt.Errorf("connection test failed: %v", err)
	}

	// Get and display connection status
	status := tester.Connector.GetStatus()
	log.Printf("📊 Connection Status:")
	log.Printf("   Connected: %v", status.IsConnected)
	log.Printf("   Database: %s", status.DatabaseName)
	log.Printf("   Connection Time: %v", status.ConnectionTime)

	// Health check
	health := tester.Connector.HealthCheck()
	log.Printf("🏥 Health Check:")
	log.Printf("   Status: %s", health.Status)
	log.Printf("   Message: %s", health.Message)
	log.Printf("   Latency: %dms", health.Latency)

	return nil
}

func runSchemaTests(tester *SchemaTester) error {
	log.Println("📁 Testing Database Schema...")

	// Test collections
	err := tester.TestCollections()
	if err != nil {
		return fmt.Errorf("collection test failed: %v", err)
	}

	// Test indexes
	err = tester.TestIndexes()
	if err != nil {
		return fmt.Errorf("index test failed: %v", err)
	}

	return nil
}

func runCRUDTests(tester *SchemaTester) error {
	log.Println("🔧 Testing CRUD Operations...")

	err := tester.TestCRUDOperations()
	if err != nil {
		return fmt.Errorf("CRUD test failed: %v", err)
	}

	return nil
}

func cleanupTestData(tester *SchemaTester) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collections := []string{
		mongoDB.CollectionUsers,
		mongoDB.CollectionWidgets,
		mongoDB.CollectionConversations,
		mongoDB.CollectionTemplates,
		mongoDB.CollectionAnalytics,
	}

	for _, collectionName := range collections {
		collection := tester.TestDB.Collection(collectionName)

		// Delete test data (you might want to be more specific in production)
		result, err := collection.DeleteMany(ctx, bson.M{
			"$or": []bson.M{
				{"email": "test@example.com"},
				{"name": bson.M{"$regex": "^Test"}},
				{"username": "testuser"},
			},
		})

		if err != nil {
			log.Printf("⚠️  Failed to clean up %s: %v", collectionName, err)
		} else {
			log.Printf("✅ Cleaned up %s: %d documents removed", collectionName, result.DeletedCount)
		}
	}
}
