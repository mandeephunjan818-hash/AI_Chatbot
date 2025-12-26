package mongoDB

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config holds MongoDB connection configuration
type Config struct {
	URI      string
	Database string
	Options  *options.ClientOptions
}

// ConnectionStatus represents the current connection state
type ConnectionStatus struct {
	IsConnected    bool
	DatabaseName   string
	ConnectionTime time.Time
	Error          string
}

// HealthCheckResult represents the health check status
type HealthCheckResult struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency_ms"`
}

// EventType represents connection event types
type EventType string

const (
	EventConnecting    EventType = "CONNECTING"
	EventConnected     EventType = "CONNECTED"
	EventDisconnecting EventType = "DISCONNECTING"
	EventDisconnected  EventType = "DISCONNECTED"
	EventError         EventType = "ERROR"
)

// ConnectionEvent represents a connection event
type ConnectionEvent struct {
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Error     string    `json:"error,omitempty"`
}

// EventHandler is a function type for handling connection events
type EventHandler func(event ConnectionEvent)

// Custom error types
type ConnectionError struct {
	Message string
}

func (e *ConnectionError) Error() string {
	return "MongoDB Connection Error: " + e.Message
}

type DisconnectionError struct {
	Message string
}

func (e *DisconnectionError) Error() string {
	return "MongoDB Disconnection Error: " + e.Message
}

type QueryError struct {
	Message string
}

func (e *QueryError) Error() string {
	return "MongoDB Query Error: " + e.Message
}
