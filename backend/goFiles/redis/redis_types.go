// redis_types.go - Redis types and interfaces
package redis

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConnector manages Redis connections
type RedisConnector struct {
	config         *RedisConfig
	client         *redis.Client
	isConnected    bool
	isConnecting   bool
	connectionTime time.Time
	mu             sync.RWMutex
	eventHandlers  []RedisEventHandler
	context        context.Context
	cancelFunc     context.CancelFunc
}

var (
	redisInstance *RedisConnector
	redisOnce     sync.Once
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Addr         string `json:"addr"`
	Password     string `json:"password"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"pool_size"`
	MinIdleConns int    `json:"min_idle_conns"`
	MaxRetries   int    `json:"max_retries"`
	DialTimeout  int    `json:"dial_timeout"`  // seconds
	ReadTimeout  int    `json:"read_timeout"`  // seconds
	WriteTimeout int    `json:"write_timeout"` // seconds
}

// RedisConnectionStatus represents the current connection state
type RedisConnectionStatus struct {
	IsConnected    bool
	Addr           string
	DB             int
	ConnectionTime time.Time
	Error          string
}

// RedisHealthCheckResult represents the health check status
type RedisHealthCheckResult struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency_ms"`
}

// RedisEventType represents connection event types
type RedisEventType string

const (
	EventConnecting    RedisEventType = "CONNECTING"
	EventConnected     RedisEventType = "CONNECTED"
	EventDisconnecting RedisEventType = "DISCONNECTING"
	EventDisconnected  RedisEventType = "DISCONNECTED"
	EventError         RedisEventType = "ERROR"
)

// RedisEvent represents a connection event
type RedisEvent struct {
	Type      RedisEventType `json:"type"`
	Timestamp time.Time      `json:"timestamp"`
	Message   string         `json:"message"`
	Error     string         `json:"error,omitempty"`
}

// RedisEventHandler is a function type for handling connection events
type RedisEventHandler func(event RedisEvent)

// Custom error types
type RedisConnectionError struct {
	Message string
}

func (e *RedisConnectionError) Error() string {
	return "Redis Connection Error: " + e.Message
}

type RedisDisconnectionError struct {
	Message string
}

func (e *RedisDisconnectionError) Error() string {
	return "Redis Disconnection Error: " + e.Message
}

type RedisOperationError struct {
	Message string
}

func (e *RedisOperationError) Error() string {
	return "Redis Operation Error: " + e.Message
}

// CacheTTL constants
const (
	TTLShort   = 5 * time.Minute    // 5 minutes
	TTLMedium  = 30 * time.Minute   // 30 minutes
	TTLLong    = 24 * time.Hour     // 24 hours
	TTLJWT     = 7 * 24 * time.Hour // 7 days
	TTLSession = 2 * time.Hour      // 2 hours
)
