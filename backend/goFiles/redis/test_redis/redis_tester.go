// redis_tester.go - Redis testing utilities
package test_redis

import (
	"context"
	"fmt"
	"log"
	"time"

	RedisConfig "backend/gofiles/redis"

	"github.com/go-redis/redis/v8"
)

// RedisTester tests Redis connections and operations
type RedisTester struct {
	Connector *RedisConfig.RedisConnector
	Client    *redis.Client
}

// NewRedisTester creates a new Redis tester
func NewRedisTester(connector *RedisConfig.RedisConnector) (*RedisTester, error) {
	client, err := connector.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis client: %v", err)
	}
	return &RedisTester{
		Connector: connector,
		Client:    client,
	}, nil
}

// TestConnection performs comprehensive connection tests
func (rt *RedisTester) TestConnection() error {
	log.Println("🔴 Testing Redis Connection...")

	// Test basic connectivity
	start := time.Now()
	status := rt.Connector.GetStatus()
	elapsed := time.Since(start)

	log.Printf("✅ Connection Status: %+v", status)
	log.Printf("⏱️  Connection latency: %v", elapsed)

	// Test ping
	start = time.Now()
	health := rt.Connector.HealthCheck()
	elapsed = time.Since(start)

	log.Printf("🏥 Health Check: %s - %s", health.Status, health.Message)
	log.Printf("⏱️  Ping latency: %v", elapsed)

	if health.Status != "healthy" {
		return fmt.Errorf("Redis health check failed: %s", health.Message)
	}

	// Get server info
	info, err := rt.Connector.GetInfo()
	if err != nil {
		log.Printf("⚠️  Failed to get Redis info: %v", err)
	} else {
		log.Printf("📊 Redis Version: %s", info["redis_version"])
		log.Printf("📊 Used Memory: %s", info["used_memory_human"])
		log.Printf("📊 Connected Clients: %s", info["connected_clients"])
	}

	log.Println("✅ Redis connection tests completed successfully")
	return nil
}

// TestBasicOperations tests basic Redis operations
func (rt *RedisTester) TestBasicOperations() error {
	log.Println("🧪 Testing Basic Redis Operations...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test SET/GET
	testKey := "test:basic:operation"
	testValue := fmt.Sprintf("test_value_%d", time.Now().Unix())

	log.Printf("  Testing SET operation...")
	err := rt.Client.Set(ctx, testKey, testValue, 60*time.Second).Err()
	if err != nil {
		return fmt.Errorf("SET operation failed: %v", err)
	}
	log.Println("    ✅ SET operation successful")

	log.Printf("  Testing GET operation...")
	retrievedValue, err := rt.Client.Get(ctx, testKey).Result()
	if err != nil {
		return fmt.Errorf("GET operation failed: %v", err)
	}

	if retrievedValue != testValue {
		return fmt.Errorf("GET value mismatch: expected %s, got %s", testValue, retrievedValue)
	}
	log.Println("    ✅ GET operation successful")

	// Test EXISTS
	log.Printf("  Testing EXISTS operation...")
	exists, err := rt.Client.Exists(ctx, testKey).Result()
	if err != nil {
		return fmt.Errorf("EXISTS operation failed: %v", err)
	}

	if exists != 1 {
		return fmt.Errorf("EXISTS returned %d, expected 1", exists)
	}
	log.Println("    ✅ EXISTS operation successful")

	// Test DEL
	log.Printf("  Testing DEL operation...")
	deleted, err := rt.Client.Del(ctx, testKey).Result()
	if err != nil {
		return fmt.Errorf("DEL operation failed: %v", err)
	}

	if deleted != 1 {
		return fmt.Errorf("DEL returned %d, expected 1", deleted)
	}
	log.Println("    ✅ DEL operation successful")

	// Verify deletion
	exists, err = rt.Client.Exists(ctx, testKey).Result()
	if err != nil {
		return fmt.Errorf("EXISTS verification failed: %v", err)
	}

	if exists != 0 {
		return fmt.Errorf("Key still exists after deletion")
	}

	log.Println("✅ Basic Redis operations test completed successfully")
	return nil
}

// TestSchemaOperations tests Redis schema operations
func (rt *RedisTester) TestSchemaOperations() error {
	log.Println("📊 Testing Redis Schema Operations...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	fmt.Println(ctx.Err())

	// Test JWT Blacklist
	log.Println("  Testing JWT Blacklist operations...")
	testToken := fmt.Sprintf("test_jwt_token_%d", time.Now().Unix())

	err := rt.Connector.SetJWTBlacklist(testToken, 5*time.Minute)
	if err != nil {
		return fmt.Errorf("JWT blacklist set failed: %v", err)
	}

	isBlacklisted, err := rt.Connector.IsJWTBlacklisted(testToken)
	if err != nil {
		return fmt.Errorf("JWT blacklist check failed: %v", err)
	}

	if !isBlacklisted {
		return fmt.Errorf("JWT token should be blacklisted")
	}
	log.Println("    ✅ JWT Blacklist operations successful")

	// Test Rate Limiting
	log.Println("  Testing Rate Limiting operations...")
	rateKey := fmt.Sprintf("test_api_key_%d", time.Now().Unix())

	// Simulate rate limit increments
	for i := 1; i <= 3; i++ {
		count, err := rt.Connector.SetRateLimit(rateKey, 10, time.Minute)
		if err != nil {
			return fmt.Errorf("Rate limit increment %d failed: %v", i, err)
		}

		current, err := rt.Connector.GetRateLimit(rateKey, time.Minute)
		if err != nil && err != redis.Nil {
			return fmt.Errorf("Rate limit get %d failed: %v", i, err)
		}

		if current != count {
			return fmt.Errorf("Rate limit count mismatch: expected %d, got %d", count, current)
		}
	}
	log.Println("    ✅ Rate Limiting operations successful")

	// Test Session Management
	log.Println("  Testing Session Management operations...")
	sessionID := fmt.Sprintf("test_session_%d", time.Now().Unix())
	sessionData := map[string]interface{}{
		"user_id":    "test_user_123",
		"email":      "test@example.com",
		"created_at": time.Now().Format(time.RFC3339),
	}

	err = rt.Connector.SessionSet(sessionID, sessionData, 30*time.Minute)
	if err != nil {
		return fmt.Errorf("Session set failed: %v", err)
	}

	retrievedSession, err := rt.Connector.SessionGet(sessionID)
	if err != nil {
		return fmt.Errorf("Session get failed: %v", err)
	}

	if retrievedSession["user_id"] != "test_user_123" {
		return fmt.Errorf("Session data mismatch")
	}

	err = rt.Connector.SessionDelete(sessionID)
	if err != nil {
		return fmt.Errorf("Session delete failed: %v", err)
	}
	log.Println("    ✅ Session Management operations successful")

	// Test Widget Cache
	log.Println("  Testing Widget Cache operations...")
	widgetID := fmt.Sprintf("test_widget_%d", time.Now().Unix())
	widgetConfig := map[string]interface{}{
		"name":          "Test Widget",
		"type":          "chat",
		"primary_color": "#3B82F6",
		"cached_at":     time.Now().Format(time.RFC3339),
	}

	err = rt.Connector.WidgetConfigCache(widgetID, widgetConfig, 10*time.Minute)
	if err != nil {
		return fmt.Errorf("Widget cache set failed: %v", err)
	}

	retrievedConfig, err := rt.Connector.WidgetConfigGet(widgetID)
	if err != nil {
		return fmt.Errorf("Widget cache get failed: %v", err)
	}

	if retrievedConfig["name"] != "Test Widget" {
		return fmt.Errorf("Widget config mismatch")
	}
	log.Println("    ✅ Widget Cache operations successful")

	// Cleanup test data
	rt.cleanupTestData()

	log.Println("✅ Redis schema operations test completed successfully")
	return nil
}

// TestPerformance tests Redis performance
func (rt *RedisTester) TestPerformance() error {
	log.Println("⚡ Testing Redis Performance...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test SET performance
	log.Println("  Testing SET performance (100 operations)...")
	start := time.Now()

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("perf:test:%d", i)
		value := fmt.Sprintf("value_%d_%d", i, time.Now().UnixNano())
		err := rt.Client.Set(ctx, key, value, 60*time.Second).Err()
		if err != nil {
			return fmt.Errorf("Performance SET failed at iteration %d: %v", i, err)
		}
	}

	elapsed := time.Since(start)
	opsPerSec := float64(100) / elapsed.Seconds()
	log.Printf("    ✅ SET performance: %v (%.2f ops/sec)", elapsed, opsPerSec)

	// Test GET performance
	log.Println("  Testing GET performance (100 operations)...")
	start = time.Now()

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("perf:test:%d", i)
		_, err := rt.Client.Get(ctx, key).Result()
		if err != nil {
			return fmt.Errorf("Performance GET failed at iteration %d: %v", i, err)
		}
	}

	elapsed = time.Since(start)
	opsPerSec = float64(100) / elapsed.Seconds()
	log.Printf("    ✅ GET performance: %v (%.2f ops/sec)", elapsed, opsPerSec)

	// Cleanup performance test data
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("perf:test:%d", i)
		rt.Client.Del(ctx, key)
	}

	log.Println("✅ Redis performance tests completed successfully")
	return nil
}

// RunAllTests runs all Redis tests
func (rt *RedisTester) RunAllTests() error {
	log.Println("🚀 Starting Comprehensive Redis Tests")
	log.Println("======================================")

	tests := []struct {
		name string
		test func() error
	}{
		{"Connection Test", rt.TestConnection},
		{"Basic Operations Test", rt.TestBasicOperations},
		{"Schema Operations Test", rt.TestSchemaOperations},
		{"Performance Test", rt.TestPerformance},
	}

	for _, t := range tests {
		log.Printf("\n📋 Running: %s", t.name)
		log.Println("--------------------------------------")

		err := t.test()
		if err != nil {
			return fmt.Errorf("❌ %s failed: %v", t.name, err)
		}

		log.Printf("✅ %s passed", t.name)
		time.Sleep(500 * time.Millisecond) // Small delay between tests
	}

	log.Println("\n======================================")
	log.Println("🎉 All Redis tests completed successfully!")
	return nil
}

// cleanupTestData removes test data
func (rt *RedisTester) cleanupTestData() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Clean up test keys
	patterns := []string{
		"test:*",
		"jwt:blacklist:*",
		"ratelimit:test_*",
		"session:test_*",
		"widget:config:test_*",
		"perf:test:*",
	}

	for _, pattern := range patterns {
		iter := rt.Client.Scan(ctx, 0, pattern, 100).Iterator()
		for iter.Next(ctx) {
			key := iter.Val()
			rt.Client.Del(ctx, key)
		}
		if err := iter.Err(); err != nil {
			log.Printf("⚠️  Failed to clean up pattern %s: %v", pattern, err)
		}
	}
}
