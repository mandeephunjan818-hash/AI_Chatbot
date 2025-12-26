// main.go - Main CLI entry point with database testing
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"backend/gofiles/config"
	"backend/gofiles/mongoDB"
	dbtest "backend/gofiles/mongoDB/tester"

	"go.mongodb.org/mongo-driver/bson"

	"backend/gofiles/redis"
	redistest "backend/gofiles/redis/test_redis"
	// "backend/gofiles/server"
)

var (
	envFile     = flag.String("env", ".env", "Environment file path")
	testOnly    = flag.Bool("run-tests", false, "Run tests only and exit")
	testType    = flag.String("test-type", "all", "Test type: mongo, redis, all")
	startServer = flag.Bool("server", false, "Start the HTTP server")
	cliMode     = flag.Bool("cli", true, "Start interactive CLI")
	verbose     = flag.Bool("run-verbose", false, "Enable verbose logging")
	cleanup     = flag.Bool("run-cleanup", false, "Clean up test data after tests")
)

func main() {
	flag.Parse()

	// Setup logging
	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("🔍 Verbose mode enabled")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	cfg.PrintConfig()

	// Initialize databases
	// log.Println("🗄️  Initializing databases...")

	// Initialize MongoDB
	mongoConnector := mongoDB.GetInstance(cfg.GetMongoDBConfig())

	// Initialize Redis
	redisConnector := redis.GetInstance(cfg.GetRedisConfig())

	// Setup signal handling for graceful shutdown
	setupSignalHandling(mongoConnector, redisConnector)

	// // Connect to databases
	// log.Println("🔗 Connecting to databases...")

	// // Connect to MongoDB
	// log.Println("  Connecting to MongoDB...")
	err = mongoConnector.Connect()
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer mongoConnector.Disconnect()

	// Connect to Redis
	log.Println("  Connecting to Redis...")
	err = redisConnector.Connect()
	if err != nil {
		log.Printf("⚠️  Failed to connect to Redis: %v", err)
		log.Println("⚠️  Continuing without Redis...")
	} else {
		defer redisConnector.Disconnect()
	}

	// Run database tests
	if *testOnly || *testType != "none" {
		err = runDatabaseTests(mongoConnector, redisConnector, *testType)
		if err != nil {
			log.Fatalf("❌ Database tests failed: %v", err)
		}

		if *testOnly {
			log.Println("✅ Tests completed successfully. Exiting.")
			return
		}
	}

	// Start HTTP server if requested
	// if *startServer {
	// 	log.Println("🚀 Starting HTTP server...")
	// 	go startHTTPServer(cfg, mongoConnector, redisConnector)
	// }

	// Start interactive CLI if requested
	if *cliMode {
		log.Println("💻 Starting interactive CLI...")
		startInteractiveCLI(mongoConnector, redisConnector)
	}

	// If neither server nor CLI started, show help
	if !*startServer && !*cliMode {
		fmt.Println("No mode selected. Use --server to start server or --cli for interactive mode.")
		flag.PrintDefaults()
	}

	// Wait for interrupt signal
	waitForInterrupt()
}

func runDatabaseTests(mongoConnector *mongoDB.Connector, redisConnector *redis.RedisConnector, testType string) error {
	log.Println("🧪 Running Database Tests...")
	log.Println("=============================")

	var mongoErr, redisErr error

	// Run MongoDB tests
	if testType == "all" || testType == "mongo" {
		log.Println("\n📊 Testing MongoDB...")
		mongoTester, err := dbtest.NewSchemaTester(mongoConnector)
		if err != nil {
			log.Printf("❌ Failed to create MongoDB tester: %v", err)
			mongoErr = err
		} else {
			err = mongoTester.RunAllTests()
			if err != nil {
				log.Printf("❌ MongoDB tests failed: %v", err)
				mongoErr = err
			} else {
				log.Println("✅ MongoDB tests passed")
			}
		}
	}

	// Run Redis tests
	if testType == "all" || testType == "redis" {
		log.Println("\n🔴 Testing Redis...")
		if redisConnector.IsConnected() {
			redisTester, err := redistest.NewRedisTester(redisConnector)
			if err != nil {
				log.Printf("❌ Failed to create Redis tester: %v", err)
				redisErr = err
			} else {
				err = redisTester.RunAllTests()
				if err != nil {
					log.Printf("❌ Redis tests failed: %v", err)
					redisErr = err
				} else {
					log.Println("✅ Redis tests passed")
				}
			}
		} else {
			log.Println("⚠️  Skipping Redis tests - not connected")
		}
	}

	// Cleanup if requested
	if *cleanup {
		log.Println("\n🧹 Cleaning up test data...")
		cleanupTestData(mongoConnector, redisConnector)
	}

	// Return combined error
	if mongoErr != nil || redisErr != nil {
		return fmt.Errorf("database tests failed: MongoDB=%v, Redis=%v", mongoErr, redisErr)
	}

	log.Println("\n🎉 All database tests completed successfully!")
	return nil
}

// func startHTTPServer(cfg *config.Config, mongoConnector *mongoDB.Connector, redisConnector *redis.RedisConnector) {
// 	// Create server instance
// 	srv := server.NewServer(cfg, mongoConnector, redisConnector)

// 	// Start server
// 	if err := srv.Start(); err != nil {
// 		log.Fatalf("❌ Failed to start server: %v", err)
// 	}
// }

func startInteractiveCLI(mongoConnector *mongoDB.Connector, redisConnector *redis.RedisConnector) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n🚀 Chat Widget System - Interactive CLI")
	fmt.Println("=========================================")
	fmt.Println("Commands:")
	fmt.Println("  status    - Show database status")
	fmt.Println("  test      - Run database tests")
	fmt.Println("  mongo     - MongoDB operations")
	fmt.Println("  redis     - Redis operations")
	fmt.Println("  schema    - Show database schema")
	fmt.Println("  stats     - Show system statistics")
	fmt.Println("  server    - Start HTTP server")
	fmt.Println("  clear     - Clear screen")
	fmt.Println("  help      - Show this help")
	fmt.Println("  exit      - Exit CLI")
	fmt.Println()

	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "status":
			showDatabaseStatus(mongoConnector, redisConnector)
		case "test":
			runDatabaseTests(mongoConnector, redisConnector, "all")
		case "mongo":
			mongoCLI(reader, mongoConnector)
		case "redis":
			redisCLI(reader, redisConnector)
		case "schema":
			showDatabaseSchema(mongoConnector)
		case "stats":
			showSystemStats(mongoConnector, redisConnector)
		// case "server":
		// 	go startHTTPServer(cfg, mongoConnector, redisConnector)
		// 	fmt.Println("✅ Server started in background")
		case "clear":
			fmt.Print("\033[H\033[2J")
		case "help":
			printHelp()
		case "exit", "quit":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Printf("❌ Unknown command: %s. Type 'help' for commands.\n", command)
		}
	}
}

func showDatabaseStatus(mongoConnector *mongoDB.Connector, redisConnector *redis.RedisConnector) {
	fmt.Println("\n📊 Database Status:")
	fmt.Println("===================")

	// MongoDB Status
	mongoStatus := mongoConnector.GetStatus()
	fmt.Println("🗄️  MongoDB:")
	fmt.Printf("  Connected: %v\n", mongoStatus.IsConnected)
	fmt.Printf("  Database: %s\n", mongoStatus.DatabaseName)
	if mongoStatus.IsConnected {
		fmt.Printf("  Since: %v\n", mongoStatus.ConnectionTime.Format("2006-01-02 15:04:05"))
	}

	mongoHealth := mongoConnector.HealthCheck()
	fmt.Printf("  Health: %s - %s\n", mongoHealth.Status, mongoHealth.Message)
	fmt.Printf("  Latency: %dms\n", mongoHealth.Latency)

	// Redis Status
	fmt.Println("\n🔴 Redis:")
	if redisConnector.IsConnected() {
		redisStatus := redisConnector.GetStatus()
		fmt.Printf("  Connected: %v\n", redisStatus.IsConnected)
		fmt.Printf("  Addr: %s\n", redisStatus.Addr)
		fmt.Printf("  DB: %d\n", redisStatus.DB)
		if redisStatus.IsConnected {
			fmt.Printf("  Since: %v\n", redisStatus.ConnectionTime.Format("2006-01-02 15:04:05"))
		}

		redisHealth := redisConnector.HealthCheck()
		fmt.Printf("  Health: %s - %s\n", redisHealth.Status, redisHealth.Message)
		fmt.Printf("  Latency: %dms\n", redisHealth.Latency)
	} else {
		fmt.Println("  Not connected")
	}
}

// listMongoCollections lists all collections in the connected database
func listMongoCollections(mongoConnector *mongoDB.Connector) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the database from the connector
	db, err := mongoConnector.GetDatabase()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	// List collection names
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		fmt.Printf("❌ Failed to list collections: %v\n", err)
		return
	}

	// Display results
	if len(collections) == 0 {
		fmt.Println("No collections found")
		return
	}

	fmt.Println("\n📦 MongoDB Collections:")
	for _, name := range collections {
		fmt.Printf("  - %s\n", name)
	}
}

// countMongoDocuments counts all documents in the specified collection
func countMongoDocuments(mongoConnector *mongoDB.Connector, collectionName string) {
	// Validate collection name
	if collectionName == "" {
		fmt.Println("❌ Collection name cannot be empty")
		fmt.Println("Usage: count <collection>")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := mongoConnector.GetDatabase()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	// Get collection reference
	collection := db.Collection(collectionName)

	// Count documents (empty filter = count all)
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		fmt.Printf("❌ Failed to count documents in '%s': %v\n", collectionName, err)
		return
	}

	fmt.Printf("📊 Collection '%s' contains %d document(s)\n", collectionName, count)
}

func mongoCLI(reader *bufio.Reader, mongoConnector *mongoDB.Connector) {
	fmt.Println("\n🗄️  MongoDB CLI")
	fmt.Println("Commands: collections, count, query, exit")

	for {
		fmt.Print("mongo> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "exit" {
			return
		}

		parts := strings.Fields(command)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "collections":
			listMongoCollections(mongoConnector)
		case "count":
			if len(parts) < 2 {
				fmt.Println("Usage: count <collection>")
				continue
			}
			countMongoDocuments(mongoConnector, parts[1])
		case "query":
			fmt.Println("Query functionality not implemented in CLI. Use the main test suite.")
		default:
			fmt.Printf("Unknown MongoDB command: %s\n", parts[0])
		}
	}
}

func redisCLI(reader *bufio.Reader, redisConnector *redis.RedisConnector) {
	if !redisConnector.IsConnected() {
		fmt.Println("❌ Redis not connected")
		return
	}

	fmt.Println("\n🔴 Redis CLI")
	fmt.Println("Commands: keys, get, set, delete, info, exit")

	for {
		fmt.Print("redis> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "exit" {
			return
		}

		parts := strings.Fields(command)
		if len(parts) == 0 {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		client, err := redisConnector.GetClient()
		if err != nil {
			fmt.Printf("❌ Redis client error: %v\n", err)
			continue
		}

		switch parts[0] {
		case "keys":
			if len(parts) < 2 {
				fmt.Println("Usage: keys <pattern>")
				continue
			}
			keys, err := client.Keys(ctx, parts[1]).Result()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Printf("Keys (%d):\n", len(keys))
				for _, key := range keys {
					fmt.Printf("  %s\n", key)
				}
			}
		case "get":
			if len(parts) < 2 {
				fmt.Println("Usage: get <key>")
				continue
			}
			val, err := client.Get(ctx, parts[1]).Result()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Printf("Value: %s\n", val)
			}
		case "set":
			if len(parts) < 3 {
				fmt.Println("Usage: set <key> <value>")
				continue
			}
			err := client.Set(ctx, parts[1], parts[2], 0).Err()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Println("OK")
			}
		case "delete":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <key>")
				continue
			}
			deleted, err := client.Del(ctx, parts[1]).Result()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Printf("Deleted %d key(s)\n", deleted)
			}
		case "info":
			info, err := redisConnector.GetInfo()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Printf("Redis Version: %s\n", info["redis_version"])
				fmt.Printf("Used Memory: %s\n", info["used_memory_human"])
				fmt.Printf("Connected Clients: %s\n", info["connected_clients"])
			}
		default:
			fmt.Printf("Unknown Redis command: %s\n", parts[0])
		}
	}
}

func showDatabaseSchema(mongoConnector *mongoDB.Connector) {
	fmt.Println("\n📁 Database Schema:")
	fmt.Println("==================")

	db, err := mongoConnector.GetDatabase()
	if err != nil {
		fmt.Printf("❌ Failed to get database: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		fmt.Printf("❌ Failed to list collections: %v\n", err)
		return
	}

	fmt.Printf("Collections (%d):\n", len(collections))
	for _, name := range collections {
		collection := db.Collection(name)
		count, _ := collection.EstimatedDocumentCount(ctx)
		fmt.Printf("  • %s (%d documents)\n", name, count)
	}

	fmt.Println("\nSchema Summary:")
	fmt.Println("  users - User accounts and profiles")
	fmt.Println("  widgets - Chat widget configurations")
	fmt.Println("  conversations - Chat conversations")
	fmt.Println("  templates - MINILMv2 templates")
	fmt.Println("  analytics_events - Usage analytics")
}

func showSystemStats(mongoConnector *mongoDB.Connector, redisConnector *redis.RedisConnector) {
	fmt.Println("\n📈 System Statistics:")
	fmt.Println("=====================")

	fmt.Println("🗄️  MongoDB Statistics:")
	db, err := mongoConnector.GetDatabase()
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		collections := []string{"users", "widgets", "conversations", "templates"}
		for _, name := range collections {
			collection := db.Collection(name)
			count, _ := collection.EstimatedDocumentCount(ctx)
			fmt.Printf("  %s: %d documents\n", name, count)
		}
	}

	fmt.Println("\n🔴 Redis Statistics:")
	if redisConnector.IsConnected() {
		info, err := redisConnector.GetInfo()
		if err == nil {
			fmt.Printf("  Memory Used: %s\n", info["used_memory_human"])
			fmt.Printf("  Connected Clients: %s\n", info["connected_clients"])
			fmt.Printf("  Total Keys: %s\n", info["db0"])
		}

		// Count keys by pattern
		patterns := []string{"jwt:*", "session:*", "widget:*", "ratelimit:*"}
		for _, pattern := range patterns {
			client, _ := redisConnector.GetClient()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			keys, _ := client.Keys(ctx, pattern).Result()
			cancel()
			fmt.Printf("  %s: %d keys\n", pattern, len(keys))
		}
	}

	fmt.Println("\n💻 System Information:")
	fmt.Printf("  Current Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("  Go Version: %s\n", runtime.Version())
}

func printHelp() {
	fmt.Println("\n📚 Available Commands:")
	fmt.Println("  status    - Show database connection status")
	fmt.Println("  test      - Run comprehensive database tests")
	fmt.Println("  mongo     - MongoDB operations CLI")
	fmt.Println("  redis     - Redis operations CLI")
	fmt.Println("  schema    - Show database schema information")
	fmt.Println("  stats     - Show system statistics")
	fmt.Println("  server    - Start HTTP server in background")
	fmt.Println("  clear     - Clear the terminal screen")
	fmt.Println("  help      - Show this help message")
	fmt.Println("  exit      - Exit the CLI")
}

func cleanupTestData(mongoConnector *mongoDB.Connector, redisConnector *redis.RedisConnector) {
	log.Println("Cleaning up test data...")

	// Clean MongoDB test data
	db, err := mongoConnector.GetDatabase()
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		collections := []string{"_health_check", "test_users", "test_widgets"}
		for _, name := range collections {
			collection := db.Collection(name)
			result, err := collection.DeleteMany(ctx, map[string]interface{}{})
			if err == nil {
				log.Printf("  MongoDB %s: deleted %d documents", name, result.DeletedCount)
			}
		}
	}

	// Clean Redis test data
	if redisConnector.IsConnected() {
		client, err := redisConnector.GetClient()
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			patterns := []string{"test:*", "jwt:blacklist:test_*", "ratelimit:test_*"}
			for _, pattern := range patterns {
				keys, _ := client.Keys(ctx, pattern).Result()
				if len(keys) > 0 {
					deleted, _ := client.Del(ctx, keys...).Result()
					log.Printf("  Redis pattern %s: deleted %d keys", pattern, deleted)
				}
			}
		}
	}

	log.Println("✅ Cleanup completed")
}

func setupSignalHandling(mongoConnector *mongoDB.Connector, redisConnector *redis.RedisConnector) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("\n⚠️  Received interrupt signal. Shutting down gracefully...")

		// Disconnect databases
		if mongoConnector != nil {
			mongoConnector.Disconnect()
		}

		if redisConnector != nil {
			redisConnector.Disconnect()
		}

		log.Println("✅ Graceful shutdown completed")
		os.Exit(0)
	}()
}

func waitForInterrupt() {
	// Create a channel to wait indefinitely
	done := make(chan bool)

	// Setup interrupt handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("\n👋 Received shutdown signal")
		done <- true
	}()

	// Wait for interrupt
	log.Println("⏳ Waiting for commands. Press Ctrl+C to exit.")
	<-done
}
