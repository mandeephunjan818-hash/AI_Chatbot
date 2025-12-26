// config.go - Configuration management with environment variables
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"backend/gofiles/mongoDB"
	"backend/gofiles/redis"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	// Server Configuration
	ServerHost         string `env:"SERVER_HOST"`
	ServerPort         int    `env:"SERVER_PORT"`
	ServerReadTimeout  int    `env:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout int    `env:"SERVER_WRITE_TIMEOUT"`
	ServerIdleTimeout  int    `env:"SERVER_IDLE_TIMEOUT"`

	// MongoDB Configuration
	MongoURI         string `env:"MONGODB_URI"`
	MongoDBName      string `env:"MONGODB_DB"`
	MongoPoolSize    int    `env:"MONGODB_POOL_SIZE"`
	MongoMaxIdleTime int    `env:"MONGODB_MAX_IDLE_TIME"`

	// Redis Configuration
	RedisAddr         string `env:"REDIS_ADDR"`
	RedisPassword     string `env:"REDIS_PASSWORD"`
	RedisDB           int    `env:"REDIS_DB"`
	RedisPoolSize     int    `env:"REDIS_POOL_SIZE"`
	RedisMinIdleConns int    `env:"REDIS_MIN_IDLE_CONNS"`

	// JWT Configuration
	JWTSecret     string `env:"JWT_SECRET"`
	JWTExpiryDays int    `env:"JWT_EXPIRY_DAYS"`

	// Feature Flags
	EnableCLI       bool `env:"ENABLE_CLI"`
	EnableWebSocket bool `env:"ENABLE_WEBSOCKET"`
	EnableTesting   bool `env:"ENABLE_TESTING"`
	EnableMetrics   bool `env:"ENABLE_METRICS"`

	// Logging
	LogLevel  string `env:"LOG_LEVEL"`
	LogFormat string `env:"LOG_FORMAT"`

	// Chat Service
	ChatServiceURL string `env:"CHAT_SERVICE_URL"`
	ChatModelPath  string `env:"CHAT_MODEL_PATH"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {

	// Try default .env file
	if err := godotenv.Load(); err == nil {
		log.Print("\n\t\t\t✅ Loaded environment from .env file\t\t\t\n\n")
	}

	config := &Config{}

	// Server Configuration
	// config.ServerHost = getEnv("SERVER_HOST", "0.0.0.0")
	// config.ServerPort = getEnvAsInt("SERVER_PORT", 8080)
	// config.ServerReadTimeout = getEnvAsInt("SERVER_READ_TIMEOUT", 30)
	// config.ServerWriteTimeout = getEnvAsInt("SERVER_WRITE_TIMEOUT", 30)
	// config.ServerIdleTimeout = getEnvAsInt("SERVER_IDLE_TIMEOUT", 60)

	// MongoDB Configuration
	config.MongoURI = getEnv("MONGODB_URI", "mongodb://localhost:27017")
	config.MongoDBName = getEnv("MONGODB_DB", "ai_chatbot")
	// config.MongoPoolSize = getEnvAsInt("MONGODB_POOL_SIZE", 20)
	// config.MongoMaxIdleTime = getEnvAsInt("MONGODB_MAX_IDLE_TIME", 30000)

	// Redis Configuration
	config.RedisAddr = getEnv("REDIS_ADDR", "localhost:6379")
	config.RedisPassword = getEnv("REDIS_PASSWORD", "")
	config.RedisDB = getEnvAsInt("REDIS_DB", 0)
	// config.RedisPoolSize = getEnvAsInt("REDIS_POOL_SIZE", 10)
	// config.RedisMinIdleConns = getEnvAsInt("REDIS_MIN_IDLE_CONNS", 3)

	// JWT Configuration
	// config.JWTSecret = getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production")
	// config.JWTExpiryDays = getEnvAsInt("JWT_EXPIRY_DAYS", 7)

	// Feature Flags
	// config.EnableCLI = getEnvAsBool("ENABLE_CLI", true)
	// config.EnableWebSocket = getEnvAsBool("ENABLE_WEBSOCKET", true)
	// config.EnableTesting = getEnvAsBool("ENABLE_TESTING", true)
	// config.EnableMetrics = getEnvAsBool("ENABLE_METRICS", true)

	// Logging
	// config.LogLevel = getEnv("LOG_LEVEL", "info")
	// config.LogFormat = getEnv("LOG_FORMAT", "text")

	// Chat Service
	// config.ChatServiceURL = getEnv("CHAT_SERVICE_URL", "http://localhost:5000")
	// config.ChatModelPath = getEnv("CHAT_MODEL_PATH", "./models")

	return config, nil
}

// GetMongoDBConfig returns MongoDB configuration
func (c *Config) GetMongoDBConfig() *mongoDB.Config {
	return &mongoDB.Config{
		URI:      c.MongoURI,
		Database: c.MongoDBName,
		Options: options.Client().
			SetConnectTimeout(30 * time.Second).
			SetMaxPoolSize(uint64(c.MongoPoolSize)).
			SetMaxConnIdleTime(time.Duration(c.MongoMaxIdleTime) * time.Millisecond).
			SetServerSelectionTimeout(30 * time.Second),
	}
}

// GetRedisConfig returns Redis configuration
func (c *Config) GetRedisConfig() *redis.RedisConfig {
	return &redis.RedisConfig{
		Addr:         c.RedisAddr,
		Password:     c.RedisPassword,
		DB:           c.RedisDB,
		PoolSize:     c.RedisPoolSize,
		MinIdleConns: c.RedisMinIdleConns,
		MaxRetries:   3,
		DialTimeout:  5,
		ReadTimeout:  3,
		WriteTimeout: 3,
	}
}

// PrintConfig prints the configuration (without secrets)
func (c *Config) PrintConfig() {
	log.Print("\n\t\t\t==============================\n\n")
	log.Print("⚙️  Application Configuration:\n\n")

	// log.Println("🌐 Server Configuration:")
	// log.Printf("  Host: %s:%d", c.ServerHost, c.ServerPort)
	// log.Printf("  Timeouts - Read: %ds, Write: %ds, Idle: %ds",
	// 	c.ServerReadTimeout, c.ServerWriteTimeout, c.ServerIdleTimeout)

	log.Println("🗄️  Database Configuration:")
	log.Printf("  MongoDB URI: %s", maskPassword(c.MongoURI))
	log.Printf("  MongoDB Database: %s", c.MongoDBName)
	// log.Printf("  MongoDB Pool Size: %d", c.MongoPoolSize)

	log.Println("🔴 Redis Configuration:")
	log.Printf("  Redis Addr: %s", c.RedisAddr)
	log.Printf("  Redis DB: %d", c.RedisDB)
	// log.Printf("  Redis Pool Size: %d", c.RedisPoolSize)
	log.Printf("  Redis Password: %s", maskString(c.RedisPassword))

	// log.Println("🔐 Security Configuration:")
	// log.Printf("  JWT Expiry: %d days", c.JWTExpiryDays)
	// log.Printf("  JWT Secret: %s", maskString(c.JWTSecret))

	// log.Println("🚀 Feature Flags:")
	// log.Printf("  CLI Enabled: %v", c.EnableCLI)
	// log.Printf("  WebSocket Enabled: %v", c.EnableWebSocket)
	// log.Printf("  Testing Enabled: %v", c.EnableTesting)
	// log.Printf("  Metrics Enabled: %v", c.EnableMetrics)

	// log.Println("📝 Logging Configuration:")
	// log.Printf("  Log Level: %s", c.LogLevel)
	// log.Printf("  Log Format: %s", c.LogFormat)

	// log.Println("💬 Chat Service:")
	// log.Printf("  Service URL: %s", c.ChatServiceURL)
	// log.Printf("  Model Path: %s", c.ChatModelPath)

	log.Print("\n\t\t\t==============================\n\n")
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func maskPassword(uri string) string {
	if strings.Contains(uri, "@") {
		parts := strings.Split(uri, "@")
		if len(parts) == 2 {
			authPart := parts[0]
			hostPart := parts[1]
			if strings.Contains(authPart, "://") {
				protocolParts := strings.Split(authPart, "://")
				if len(protocolParts) == 2 {
					protocol := protocolParts[0]
					credentials := protocolParts[1]
					if strings.Contains(credentials, ":") {
						credParts := strings.Split(credentials, ":")
						if len(credParts) >= 2 {
							username := credParts[0]
							return fmt.Sprintf("%s://%s:****@%s", protocol, username, hostPart)
						}
					}
				}
			}
		}
	}
	return uri
}

func maskString(value string) string {
	if value == "" {
		return "(empty)"
	}
	if len(value) > 8 {
		return value[:3] + "****" + value[len(value)-3:]
	}
	return "****"
}
