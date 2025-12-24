// schema.go
package main

import "time"

// Conversation represents a chat session document in MongoDB
type Conversation struct {
	AppID       string                 `bson:"app_id"`
	UserToken   string                 `bson:"user_token"`
	SessionID   string                 `bson:"session_id"`
	UserMessage string                 `bson:"user_message"`
	BotReply    string                 `bson:"bot_reply"`
	Metadata    map[string]interface{} `bson:"metadata"`
	Timestamp   time.Time              `bson:"timestamp"`
}
