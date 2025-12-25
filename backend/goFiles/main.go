package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Parse command line flags
	runTestOnly := flag.Bool("test", false, "Run database tests only and exit")
	quickTest := flag.Bool("quick", false, "Run quick connection test only")
	skipTest := flag.Bool("skip-test", false, "Skip initial connection test (not recommended)")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *help {
		printHelp()
		return
	}

	// Initialize MongoDB connection
	log.Println("🔌 Initializing MongoDB connection...")
	InitMongoDB()
	defer func() {
		if err := MongoDb.Client().Disconnect(context.Background()); err != nil {
			log.Println("⚠️  Failed to disconnect from MongoDB:", err)
		}
		log.Println("👋 MongoDB connection closed")
	}()

	// Handle test-only mode
	if *runTestOnly {
		log.Println("🧪 Running full database test suite...")
		if RunFullTest() {
			os.Exit(0) // Exit with success
		} else {
			os.Exit(1) // Exit with error
		}
	}

	if *quickTest {
		log.Println("⚡ Running quick connection test...")
		if QuickTest() {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	// Always test connection before starting server (unless skipped)
	if !*skipTest {
		log.Println("🔍 Testing database connection before starting server...")
		if !QuickTest() {
			log.Fatal("❌ Database connection test failed. Server not started.")
		}
		log.Println("✅ Database connection verified. Starting server...")
	} else {
		log.Println("⚠️  Skipping initial connection test (not recommended)")
	}

	// Start the HTTP server
	StartServer()
}

func printHelp() {
	fmt.Println("AI Chatbot API Server")
	fmt.Println("======================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run . [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -test        Run full database test and exit")
	fmt.Println("  -quick       Run quick connection test and exit")
	fmt.Println("  -skip-test   Skip initial connection test (not recommended)")
	fmt.Println("  -help        Show this help message")
	fmt.Println()
	fmt.Println("HTTP Endpoints (when server is running):")
	fmt.Println("  GET  /health     - Basic health check")
	fmt.Println("  GET  /test-db    - Quick database connection test")
	fmt.Println("  GET  /db-info    - Database connection information")
	fmt.Println("  GET  /full-test  - Run full database test")
	fmt.Println("  POST /test-log   - Log a test conversation")
}
