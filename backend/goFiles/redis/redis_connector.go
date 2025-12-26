// redis_connector.go - Redis connection management
package redis

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// GetInstance returns a singleton instance of Redis Connector
func GetInstance(config *RedisConfig) *RedisConnector {
	redisOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		redisInstance = &RedisConnector{
			config:        config,
			isConnected:   false,
			context:       ctx,
			cancelFunc:    cancel,
			eventHandlers: make([]RedisEventHandler, 0),
		}
	})
	return redisInstance
}

// Connect establishes connection to Redis
func (rc *RedisConnector) Connect() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.isConnected {
		rc.emitEvent(RedisEvent{
			Type:      EventConnected,
			Timestamp: time.Now(),
			Message:   "Already connected to Redis",
		})
		return nil
	}

	if rc.isConnecting {
		return &RedisConnectionError{Message: "Redis connection already in progress"}
	}

	rc.isConnecting = true
	rc.emitEvent(RedisEvent{
		Type:      EventConnecting,
		Timestamp: time.Now(),
		Message:   "Connecting to Redis...",
	})

	// Create Redis options
	options := &redis.Options{
		Addr:     rc.config.Addr,
		Password: rc.config.Password,
		DB:       rc.config.DB,
	}

	if rc.config.MaxRetries > 0 {
		options.MaxRetries = rc.config.MaxRetries
	}

	if rc.config.PoolSize > 0 {
		options.PoolSize = rc.config.PoolSize
	}

	if rc.config.MinIdleConns > 0 {
		options.MinIdleConns = rc.config.MinIdleConns
	}

	if rc.config.DialTimeout > 0 {
		options.DialTimeout = time.Duration(rc.config.DialTimeout) * time.Second
	}

	if rc.config.ReadTimeout > 0 {
		options.ReadTimeout = time.Duration(rc.config.ReadTimeout) * time.Second
	}

	if rc.config.WriteTimeout > 0 {
		options.WriteTimeout = time.Duration(rc.config.WriteTimeout) * time.Second
	}

	// Create Redis client
	client := redis.NewClient(options)

	// Test connection
	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		rc.isConnecting = false
		rc.emitEvent(RedisEvent{
			Type:      EventError,
			Timestamp: time.Now(),
			Message:   "Failed to connect to Redis",
			Error:     err.Error(),
		})
		return &RedisConnectionError{Message: fmt.Sprintf("Failed to connect: %v", err)}
	}

	rc.client = client
	rc.isConnected = true
	rc.isConnecting = false
	rc.connectionTime = time.Now()

	rc.emitEvent(RedisEvent{
		Type:      EventConnected,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("✅ Successfully connected to Redis at: %s", rc.config.Addr),
	})

	log.Printf("✅ Redis connected successfully to: %s (DB: %d)", rc.config.Addr, rc.config.DB)
	return nil
}

// Disconnect closes the Redis connection
func (rc *RedisConnector) Disconnect() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if !rc.isConnected {
		rc.emitEvent(RedisEvent{
			Type:      EventDisconnected,
			Timestamp: time.Now(),
			Message:   "Already disconnected from Redis",
		})
		return nil
	}

	rc.emitEvent(RedisEvent{
		Type:      EventDisconnecting,
		Timestamp: time.Now(),
		Message:   "Disconnecting from Redis...",
	})

	err := rc.client.Close()
	if err != nil {
		rc.emitEvent(RedisEvent{
			Type:      EventError,
			Timestamp: time.Now(),
			Message:   "Failed to disconnect from Redis",
			Error:     err.Error(),
		})
		return &RedisDisconnectionError{Message: fmt.Sprintf("Failed to disconnect: %v", err)}
	}

	rc.isConnected = false
	rc.client = nil

	rc.emitEvent(RedisEvent{
		Type:      EventDisconnected,
		Timestamp: time.Now(),
		Message:   "✅ Successfully disconnected from Redis",
	})

	log.Println("✅ Redis disconnected successfully")
	return nil
}

// GetClient returns the Redis client instance
func (rc *RedisConnector) GetClient() (*redis.Client, error) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	if !rc.isConnected || rc.client == nil {
		return nil, &RedisConnectionError{Message: "Redis client not connected. Call Connect() first."}
	}

	return rc.client, nil
}

// HealthCheck performs a health check on the connection
func (rc *RedisConnector) HealthCheck() RedisHealthCheckResult {
	start := time.Now()

	if !rc.isConnected {
		return RedisHealthCheckResult{
			Status:    "unhealthy",
			Message:   "Not connected to Redis",
			Timestamp: time.Now(),
			Latency:   0,
		}
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	_, err := rc.client.Ping(ctx).Result()
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return RedisHealthCheckResult{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("Redis ping failed: %v", err),
			Timestamp: time.Now(),
			Latency:   latency,
		}
	}

	return RedisHealthCheckResult{
		Status:    "healthy",
		Message:   "Redis connection is healthy",
		Timestamp: time.Now(),
		Latency:   latency,
	}
}

// GetStatus returns the current connection status
func (rc *RedisConnector) GetStatus() RedisConnectionStatus {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	status := RedisConnectionStatus{
		IsConnected: rc.isConnected,
		Addr:        rc.config.Addr,
		DB:          rc.config.DB,
	}

	if rc.isConnected {
		status.ConnectionTime = rc.connectionTime
	}

	return status
}

// IsConnected returns true if connected to Redis
func (rc *RedisConnector) IsConnected() bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.isConnected
}

// AddEventHandler adds an event handler for connection events
func (rc *RedisConnector) AddEventHandler(handler RedisEventHandler) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.eventHandlers = append(rc.eventHandlers, handler)
}

// emitEvent sends events to all registered handlers
func (rc *RedisConnector) emitEvent(event RedisEvent) {
	for _, handler := range rc.eventHandlers {
		handler(event)
	}
}

// Close cleans up resources
func (rc *RedisConnector) Close() error {
	if rc.cancelFunc != nil {
		rc.cancelFunc()
	}

	if rc.isConnected {
		return rc.Disconnect()
	}
	return nil
}

// Common Redis operations with our schema

// SetJWTBlacklist adds a JWT token to blacklist
func (rc *RedisConnector) SetJWTBlacklist(token string, expiry time.Duration) error {
	client, err := rc.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("jwt:blacklist:%s", token)
	return client.Set(ctx, key, "1", expiry).Err()
}

// IsJWTBlacklisted checks if a JWT token is blacklisted
func (rc *RedisConnector) IsJWTBlacklisted(token string) (bool, error) {
	client, err := rc.GetClient()
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("jwt:blacklist:%s", token)
	exists, err := client.Exists(ctx, key).Result()
	return exists == 1, err
}

// SetRateLimit sets rate limit for a key
func (rc *RedisConnector) SetRateLimit(key string, limit int, window time.Duration) (int64, error) {
	client, err := rc.GetClient()
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	redisKey := fmt.Sprintf("ratelimit:%s:%d", key, time.Now().Unix()/int64(window.Seconds()))
	return client.Incr(ctx, redisKey).Result()
}

// GetRateLimit gets current rate limit count
func (rc *RedisConnector) GetRateLimit(key string, window time.Duration) (int64, error) {
	client, err := rc.GetClient()
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	redisKey := fmt.Sprintf("ratelimit:%s:%d", key, time.Now().Unix()/int64(window.Seconds()))
	return client.Get(ctx, redisKey).Int64()
}

// CacheSet sets a value in cache
func (rc *RedisConnector) CacheSet(key string, value interface{}, expiry time.Duration) error {
	client, err := rc.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	return client.Set(ctx, key, value, expiry).Err()
}

// CacheGet gets a value from cache
func (rc *RedisConnector) CacheGet(key string) (string, error) {
	client, err := rc.GetClient()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	return client.Get(ctx, key).Result()
}

// CacheDelete deletes a value from cache
func (rc *RedisConnector) CacheDelete(key string) error {
	client, err := rc.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	return client.Del(ctx, key).Err()
}

// CacheExists checks if a key exists in cache
func (rc *RedisConnector) CacheExists(key string) (bool, error) {
	client, err := rc.GetClient()
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	exists, err := client.Exists(ctx, key).Result()
	return exists == 1, err
}

// SessionSet sets a user session
func (rc *RedisConnector) SessionSet(sessionID string, userData map[string]interface{}, expiry time.Duration) error {
	client, err := rc.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("session:%s", sessionID)
	return client.HSet(ctx, key, userData).Err()
}

// SessionGet gets a user session
func (rc *RedisConnector) SessionGet(sessionID string) (map[string]string, error) {
	client, err := rc.GetClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("session:%s", sessionID)
	return client.HGetAll(ctx, key).Result()
}

// SessionDelete deletes a user session
func (rc *RedisConnector) SessionDelete(sessionID string) error {
	client, err := rc.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("session:%s", sessionID)
	return client.Del(ctx, key).Err()
}

// WidgetConfigCache caches widget configuration
func (rc *RedisConnector) WidgetConfigCache(widgetID string, config map[string]interface{}, expiry time.Duration) error {
	client, err := rc.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("widget:config:%s", widgetID)
	return client.HSet(ctx, key, config).Err()
}

// WidgetConfigGet gets cached widget configuration
func (rc *RedisConnector) WidgetConfigGet(widgetID string) (map[string]string, error) {
	client, err := rc.GetClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("widget:config:%s", widgetID)
	return client.HGetAll(ctx, key).Result()
}

// ConversationCache caches conversation state
func (rc *RedisConnector) ConversationCache(conversationID string, state map[string]interface{}, expiry time.Duration) error {
	client, err := rc.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("conversation:state:%s", conversationID)
	return client.HSet(ctx, key, state).Err()
}

// ConversationGet gets cached conversation state
func (rc *RedisConnector) ConversationGet(conversationID string) (map[string]string, error) {
	client, err := rc.GetClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("conversation:state:%s", conversationID)
	return client.HGetAll(ctx, key).Result()
}

// GetInfo returns Redis server information
func (rc *RedisConnector) GetInfo() (map[string]string, error) {
	client, err := rc.GetClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(rc.context, 5*time.Second)
	defer cancel()

	info, err := client.Info(ctx).Result()
	if err != nil {
		return nil, err
	}

	// Parse info into map
	result := make(map[string]string)
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				result[key] = value
			}
		}
	}

	return result, nil
}
