package mongoDB

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connector manages MongoDB connections
type Connector struct {
	config         *Config
	client         *mongo.Client
	database       *mongo.Database
	isConnected    bool
	isConnecting   bool
	connectionTime time.Time
	mu             sync.RWMutex
	eventHandlers  []EventHandler
	context        context.Context
	cancelFunc     context.CancelFunc
}

var (
	instance *Connector
	once     sync.Once
)

// GetInstance returns a singleton instance of MongoDB Connector
func GetInstance(config *Config) *Connector {
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		instance = &Connector{
			config:      config,
			isConnected: false,
			context:     ctx,
			cancelFunc:  cancel,
		}
	})
	return instance
}

// Connect establishes connection to MongoDB
func (c *Connector) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		c.emitEvent(ConnectionEvent{
			Type:      EventConnected,
			Timestamp: time.Now(),
			Message:   "Already connected to MongoDB",
		})
		return nil
	}

	if c.isConnecting {
		return &ConnectionError{Message: "Connection already in progress"}
	}

	c.isConnecting = true
	c.emitEvent(ConnectionEvent{
		Type:      EventConnecting,
		Timestamp: time.Now(),
		Message:   "Connecting to MongoDB...",
	})

	// Create MongoDB client options
	clientOptions := options.Client()
	if c.config.Options != nil {
		clientOptions = c.config.Options
	}

	// Apply URI
	if c.config.URI != "" {
		clientOptions.ApplyURI(c.config.URI)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(c.context, 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		c.isConnecting = false
		c.emitEvent(ConnectionEvent{
			Type:      EventError,
			Timestamp: time.Now(),
			Message:   "Failed to connect to MongoDB",
			Error:     err.Error(),
		})
		return &ConnectionError{Message: fmt.Sprintf("Failed to connect: %v", err)}
	}

	// Verify connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		c.isConnecting = false
		c.emitEvent(ConnectionEvent{
			Type:      EventError,
			Timestamp: time.Now(),
			Message:   "Failed to ping MongoDB",
			Error:     err.Error(),
		})
		return &ConnectionError{Message: fmt.Sprintf("Failed to ping: %v", err)}
	}

	c.client = client
	c.database = client.Database(c.config.Database)
	c.isConnected = true
	c.isConnecting = false
	c.connectionTime = time.Now()

	c.emitEvent(ConnectionEvent{
		Type:      EventConnected,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("✅ Successfully connected to database: %s", c.config.Database),
	})

	log.Printf("✅ MongoDB connected successfully to database: %s", c.config.Database)
	return nil
}

// Disconnect closes the MongoDB connection
func (c *Connector) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		c.emitEvent(ConnectionEvent{
			Type:      EventDisconnected,
			Timestamp: time.Now(),
			Message:   "Already disconnected from MongoDB",
		})
		return nil
	}

	c.emitEvent(ConnectionEvent{
		Type:      EventDisconnecting,
		Timestamp: time.Now(),
		Message:   "Disconnecting from MongoDB...",
	})

	ctx, cancel := context.WithTimeout(c.context, 5*time.Second)
	defer cancel()

	err := c.client.Disconnect(ctx)
	if err != nil {
		c.emitEvent(ConnectionEvent{
			Type:      EventError,
			Timestamp: time.Now(),
			Message:   "Failed to disconnect from MongoDB",
			Error:     err.Error(),
		})
		return &DisconnectionError{Message: fmt.Sprintf("Failed to disconnect: %v", err)}
	}

	c.isConnected = false
	c.client = nil
	c.database = nil

	c.emitEvent(ConnectionEvent{
		Type:      EventDisconnected,
		Timestamp: time.Now(),
		Message:   "✅ Successfully disconnected from MongoDB",
	})

	log.Println("✅ MongoDB disconnected successfully")
	return nil
}

// GetDatabase returns the database instance
func (c *Connector) GetDatabase() (*mongo.Database, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.isConnected || c.database == nil {
		return nil, &ConnectionError{Message: "Database not connected. Call Connect() first."}
	}

	return c.database, nil
}

// GetClient returns the MongoDB client instance
func (c *Connector) GetClient() (*mongo.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.isConnected || c.client == nil {
		return nil, &ConnectionError{Message: "Client not connected. Call Connect() first."}
	}

	return c.client, nil
}

// GetStatus returns the current connection status
func (c *Connector) GetStatus() ConnectionStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status := ConnectionStatus{
		IsConnected:  c.isConnected,
		DatabaseName: c.config.Database,
	}

	if c.isConnected {
		status.ConnectionTime = c.connectionTime
	}

	return status
}

// HealthCheck performs a health check on the connection
func (c *Connector) HealthCheck() HealthCheckResult {
	start := time.Now()

	if !c.isConnected {
		return HealthCheckResult{
			Status:    "unhealthy",
			Message:   "Not connected to MongoDB",
			Timestamp: time.Now(),
			Latency:   0,
		}
	}

	ctx, cancel := context.WithTimeout(c.context, 5*time.Second)
	defer cancel()

	err := c.client.Ping(ctx, readpref.Primary())
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return HealthCheckResult{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("MongoDB ping failed: %v", err),
			Timestamp: time.Now(),
			Latency:   latency,
		}
	}

	return HealthCheckResult{
		Status:    "healthy",
		Message:   "MongoDB connection is healthy",
		Timestamp: time.Now(),
		Latency:   latency,
	}
}

// IsConnected returns true if connected to MongoDB
func (c *Connector) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// AddEventHandler adds an event handler for connection events
func (c *Connector) AddEventHandler(handler EventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.eventHandlers = append(c.eventHandlers, handler)
}

// emitEvent sends events to all registered handlers
func (c *Connector) emitEvent(event ConnectionEvent) {
	for _, handler := range c.eventHandlers {
		handler(event)
	}
}

// Close cleans up resources
func (c *Connector) Close() error {
	if c.cancelFunc != nil {
		c.cancelFunc()
	}

	if c.isConnected {
		return c.Disconnect()
	}
	return nil
}
