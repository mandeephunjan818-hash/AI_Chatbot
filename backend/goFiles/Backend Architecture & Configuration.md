# Backend Architecture & Configuration

This section details the Go backend service, specifically regarding dependencies, environment configuration, and database schema.

## Dependencies

The application leverages the following key Go modules:

- **go.mongodb.org/mongo-driver**: Official MongoDB driver.

- **github.com/joho/godotenv**: Loads environment variables from `.env`.

- **github.com/redis/go-redis**: Redis client for caching/session management.

## Environment Configuration

The application requires a `.env` file in the root directory.

```bash
# Database Configuration

# The connection string for the MongoDB instance
MONGODB_URI=mongodb://localhost:27017

# The specific database name to use within the cluster
MONGODB_DB=ai_chatbot
```

## Database Schema

The primary data model for logging chat sessions in MongoDB is the `Conversation` struct.

```go
package main

import "time"

// Conversation represents a chat session document
type Conversation struct {
    AppID       string                 `bson:"app_id"`
    UserToken   string                 `bson:"user_token"`
    SessionID   string                 `bson:"session_id"`
    UserMessage string                 `bson:"user_message"`
    BotReply    string                 `bson:"bot_reply"`
    Metadata    map[string]interface{} `bson:"metadata"`
    Timestamp   time.Time              `bson:"timestamp"`
}
```

