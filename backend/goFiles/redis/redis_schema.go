// redis_schema.go - Redis data schema definitions
package redis

import (
	"fmt"
	"time"
)

// Key patterns for Redis data organization
const (
	// JWT Management
	KeyJWTBlacklist = "jwt:blacklist:%s" // token

	// Rate Limiting
	KeyRateLimit = "ratelimit:%s:%d" // identifier:timestamp

	// User Sessions
	KeyUserSession = "session:%s"     // session_id
	KeyUserTokens  = "user:tokens:%s" // user_id

	// Widget Cache
	KeyWidgetConfig = "widget:config:%s" // widget_id
	KeyWidgetStats  = "widget:stats:%s"  // widget_id
	KeyWidgetActive = "widget:active"    // sorted set

	// Conversation Cache
	KeyConversationState = "conversation:state:%s" // conversation_id
	KeyConversationLock  = "conversation:lock:%s"  // conversation_id

	// API Keys
	KeyAPIKeyValid     = "apikey:valid:%s"        // api_key
	KeyAPIKeyRateLimit = "apikey:ratelimit:%s:%d" // api_key:timestamp

	// Templates
	KeyTemplateCache   = "template:%s"      // template_id
	KeyTemplatePopular = "template:popular" // sorted set

	// Analytics
	KeyAnalyticsEvents   = "analytics:events:%s" // date (YYYY-MM-DD)
	KeyAnalyticsRealTime = "analytics:realtime"  // stream

	// System
	KeySystemHealth  = "system:health"
	KeySystemMetrics = "system:metrics"
	KeySystemAlerts  = "system:alerts"

	// Queue Names
	QueueEmail        = "queue:email"
	QueueNotification = "queue:notification"
	QueueAnalytics    = "queue:analytics"

	// Lock Keys (for distributed locks)
	LockWidgetUpdate = "lock:widget:update:%s" // widget_id
	LockUserUpdate   = "lock:user:update:%s"   // user_id
)

// Data structures for different Redis use cases
type JWTCache struct {
	Token       string    `json:"token"`
	UserID      string    `json:"user_id"`
	ExpiresAt   time.Time `json:"expires_at"`
	Blacklisted bool      `json:"blacklisted"`
}

type SessionData struct {
	UserID     string                 `json:"user_id"`
	Email      string                 `json:"email"`
	Username   string                 `json:"username"`
	Roles      []string               `json:"roles"`
	CreatedAt  time.Time              `json:"created_at"`
	LastAccess time.Time              `json:"last_access"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type RateLimitData struct {
	Key       string    `json:"key"`
	Count     int64     `json:"count"`
	Limit     int64     `json:"limit"`
	Window    int64     `json:"window"` // seconds
	ResetAt   time.Time `json:"reset_at"`
	Remaining int64     `json:"remaining"`
}

type WidgetCache struct {
	WidgetID   string                 `json:"widget_id"`
	Config     map[string]interface{} `json:"config"`
	Appearance map[string]interface{} `json:"appearance"`
	Behavior   map[string]interface{} `json:"behavior"`
	CachedAt   time.Time              `json:"cached_at"`
	ExpiresAt  time.Time              `json:"expires_at"`
}

type ConversationCache struct {
	ConversationID string                 `json:"conversation_id"`
	WidgetID       string                 `json:"widget_id"`
	Messages       []MessageCache         `json:"messages"`
	State          map[string]interface{} `json:"state"`
	LastUpdated    time.Time              `json:"last_updated"`
	ExpiresAt      time.Time              `json:"expires_at"`
}

type MessageCache struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
	Read      bool      `json:"read"`
}

type TemplateCache struct {
	TemplateID string    `json:"template_id"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	Variables  []string  `json:"variables"`
	UsageCount int64     `json:"usage_count"`
	CachedAt   time.Time `json:"cached_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type AnalyticsEvent struct {
	EventType string                 `json:"event_type"`
	WidgetID  string                 `json:"widget_id"`
	UserID    string                 `json:"user_id"`
	VisitorID string                 `json:"visitor_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Redis usage statistics
type RedisStats struct {
	UsedMemory        int64     `json:"used_memory"`
	UsedMemoryHuman   string    `json:"used_memory_human"`
	TotalKeys         int64     `json:"total_keys"`
	Uptime            int64     `json:"uptime"`
	ConnectedClients  int64     `json:"connected_clients"`
	CommandsProcessed int64     `json:"commands_processed"`
	Timestamp         time.Time `json:"timestamp"`
}

// GetKeyPattern returns the formatted Redis key
func GetKeyPattern(pattern string, args ...interface{}) string {
	return fmt.Sprintf(pattern, args...)
}
