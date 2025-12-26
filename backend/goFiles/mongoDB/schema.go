// schema.go - Database schemas for the architecture
package mongoDB

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User Schema
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	Username     string             `bson:"username" json:"username"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	FullName     string             `bson:"full_name,omitempty" json:"full_name,omitempty"`
	Company      string             `bson:"company,omitempty" json:"company,omitempty"`
	Plan         string             `bson:"plan" json:"plan"` // free, pro, enterprise
	APIKeys      []APIKey           `bson:"api_keys,omitempty" json:"api_keys,omitempty"`
	Settings     UserSettings       `bson:"settings" json:"settings"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	LastLogin    *time.Time         `bson:"last_login,omitempty" json:"last_login,omitempty"`
	IsActive     bool               `bson:"is_active" json:"is_active"`
	IsVerified   bool               `bson:"is_verified" json:"is_verified"`
}

type APIKey struct {
	Key       string    `bson:"key" json:"key"`
	Name      string    `bson:"name" json:"name"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	LastUsed  time.Time `bson:"last_used,omitempty" json:"last_used,omitempty"`
	IsActive  bool      `bson:"is_active" json:"is_active"`
	RateLimit int       `bson:"rate_limit" json:"rate_limit"` // requests per minute
}

type UserSettings struct {
	Theme              string `bson:"theme" json:"theme"`
	EmailNotifications bool   `bson:"email_notifications" json:"email_notifications"`
	MaxWidgets         int    `bson:"max_widgets" json:"max_widgets"`
	StorageLimit       int64  `bson:"storage_limit" json:"storage_limit"` // in MB
}

// Widget Schema
type Widget struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name        string             `bson:"name" json:"name"`
	Type        string             `bson:"type" json:"type"` // chat, feedback, support, custom
	Description string             `bson:"description,omitempty" json:"description,omitempty"`

	// Configuration
	Config     WidgetConfig     `bson:"config" json:"config"`
	Appearance AppearanceConfig `bson:"appearance" json:"appearance"`
	Behavior   BehaviorConfig   `bson:"behavior" json:"behavior"`

	// Placement
	Domains   []string `bson:"domains" json:"domains"`
	Pages     []string `bson:"pages,omitempty" json:"pages,omitempty"` // specific page paths
	Placement string   `bson:"placement" json:"placement"`             // bottom-right, bottom-left, custom

	// Status
	IsActive    bool `bson:"is_active" json:"is_active"`
	IsPublished bool `bson:"is_published" json:"is_published"`

	// Analytics
	Stats WidgetStats `bson:"stats,omitempty" json:"stats,omitempty"`

	// Timestamps
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`
	PublishedAt *time.Time `bson:"published_at,omitempty" json:"published_at,omitempty"`
}

type WidgetConfig struct {
	WelcomeMessage   string `bson:"welcome_message" json:"welcome_message"`
	InitialPrompt    string `bson:"initial_prompt,omitempty" json:"initial_prompt,omitempty"`
	ResponseDelay    int    `bson:"response_delay" json:"response_delay"` // milliseconds
	MaxMessageLength int    `bson:"max_message_length" json:"max_message_length"`
	Language         string `bson:"language" json:"language"`
	Timezone         string `bson:"timezone" json:"timezone"`

	// AI Settings
	Model       string  `bson:"model" json:"model"` // MINILMv2, other models
	Temperature float64 `bson:"temperature" json:"temperature"`
	MaxTokens   int     `bson:"max_tokens" json:"max_tokens"`

	// Business Rules
	BusinessHours BusinessHours `bson:"business_hours,omitempty" json:"business_hours,omitempty"`
	Holidays      []time.Time   `bson:"holidays,omitempty" json:"holidays,omitempty"`

	// Integrations
	WebhookURL      string `bson:"webhook_url,omitempty" json:"webhook_url,omitempty"`
	EmailForwarding string `bson:"email_forwarding,omitempty" json:"email_forwarding,omitempty"`
}

type AppearanceConfig struct {
	PrimaryColor   string `bson:"primary_color" json:"primary_color"`
	SecondaryColor string `bson:"secondary_color" json:"secondary_color"`
	TextColor      string `bson:"text_color" json:"text_color"`
	ButtonText     string `bson:"button_text" json:"button_text"`
	Icon           string `bson:"icon" json:"icon"` // emoji or image URL
	Position       string `bson:"position" json:"position"`
	Animation      string `bson:"animation,omitempty" json:"animation,omitempty"`
	CustomCSS      string `bson:"custom_css,omitempty" json:"custom_css,omitempty"`
}

type BehaviorConfig struct {
	AutoOpen         bool `bson:"auto_open" json:"auto_open"`
	DelaySeconds     int  `bson:"delay_seconds" json:"delay_seconds"`
	ShowOnMobile     bool `bson:"show_on_mobile" json:"show_on_mobile"`
	ShowOnTablet     bool `bson:"show_on_tablet" json:"show_on_tablet"`
	ShowOnDesktop    bool `bson:"show_on_desktop" json:"show_on_desktop"`
	ExitIntent       bool `bson:"exit_intent" json:"exit_intent"`
	PageScroll       bool `bson:"page_scroll" json:"page_scroll"`
	ScrollPercentage int  `bson:"scroll_percentage" json:"scroll_percentage"`
}

type BusinessHours struct {
	Enabled        bool          `bson:"enabled" json:"enabled"`
	Timezone       string        `bson:"timezone" json:"timezone"`
	Schedule       []DaySchedule `bson:"schedule" json:"schedule"`
	OfflineMessage string        `bson:"offline_message" json:"offline_message"`
}

type DaySchedule struct {
	Day       string `bson:"day" json:"day"`               // Monday, Tuesday, etc.
	OpenTime  string `bson:"open_time" json:"open_time"`   // "09:00"
	CloseTime string `bson:"close_time" json:"close_time"` // "17:00"
	IsOpen    bool   `bson:"is_open" json:"is_open"`
}

type WidgetStats struct {
	TotalConversations  int        `bson:"total_conversations" json:"total_conversations"`
	ActiveConversations int        `bson:"active_conversations" json:"active_conversations"`
	MessagesSent        int        `bson:"messages_sent" json:"messages_sent"`
	MessagesReceived    int        `bson:"messages_received" json:"messages_received"`
	AvgResponseTime     float64    `bson:"avg_response_time" json:"avg_response_time"`
	SatisfactionRate    float64    `bson:"satisfaction_rate,omitempty" json:"satisfaction_rate,omitempty"`
	LastActive          *time.Time `bson:"last_active,omitempty" json:"last_active,omitempty"`
}

// Conversation Schema
type Conversation struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	WidgetID primitive.ObjectID `bson:"widget_id" json:"widget_id"`
	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`

	// Visitor Information
	VisitorID string          `bson:"visitor_id" json:"visitor_id"`
	VisitorIP string          `bson:"visitor_ip,omitempty" json:"visitor_ip,omitempty"`
	UserAgent string          `bson:"user_agent,omitempty" json:"user_agent,omitempty"`
	PageURL   string          `bson:"page_url" json:"page_url"`
	Referrer  string          `bson:"referrer,omitempty" json:"referrer,omitempty"`
	Location  VisitorLocation `bson:"location,omitempty" json:"location,omitempty"`
	Device    DeviceInfo      `bson:"device,omitempty" json:"device,omitempty"`

	// Conversation Data
	Messages []Message `bson:"messages" json:"messages"`
	Status   string    `bson:"status" json:"status"` // active, resolved, archived, spam
	Category string    `bson:"category,omitempty" json:"category,omitempty"`
	Tags     []string  `bson:"tags,omitempty" json:"tags,omitempty"`
	Priority string    `bson:"priority,omitempty" json:"priority,omitempty"` // low, medium, high, urgent

	// AI Processing
	Intent     string  `bson:"intent,omitempty" json:"intent,omitempty"`
	Sentiment  string  `bson:"sentiment,omitempty" json:"sentiment,omitempty"`
	Confidence float64 `bson:"confidence,omitempty" json:"confidence,omitempty"`
	ModelUsed  string  `bson:"model_used,omitempty" json:"model_used,omitempty"`

	// Human Takeover
	AssignedTo      *primitive.ObjectID `bson:"assigned_to,omitempty" json:"assigned_to,omitempty"`
	IsHumanTakeover bool                `bson:"is_human_takeover" json:"is_human_takeover"`
	TakeoverReason  string              `bson:"takeover_reason,omitempty" json:"takeover_reason,omitempty"`

	// Feedback
	Rating   int    `bson:"rating,omitempty" json:"rating,omitempty"` // 1-5
	Feedback string `bson:"feedback,omitempty" json:"feedback,omitempty"`

	// Timestamps
	CreatedAt       time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `bson:"updated_at" json:"updated_at"`
	ResolvedAt      *time.Time `bson:"resolved_at,omitempty" json:"resolved_at,omitempty"`
	FirstResponseAt *time.Time `bson:"first_response_at,omitempty" json:"first_response_at,omitempty"`
}

type VisitorLocation struct {
	Country   string  `bson:"country,omitempty" json:"country,omitempty"`
	Region    string  `bson:"region,omitempty" json:"region,omitempty"`
	City      string  `bson:"city,omitempty" json:"city,omitempty"`
	Latitude  float64 `bson:"latitude,omitempty" json:"latitude,omitempty"`
	Longitude float64 `bson:"longitude,omitempty" json:"longitude,omitempty"`
	Timezone  string  `bson:"timezone,omitempty" json:"timezone,omitempty"`
}

type DeviceInfo struct {
	Type       string `bson:"type,omitempty" json:"type,omitempty"` // mobile, tablet, desktop
	OS         string `bson:"os,omitempty" json:"os,omitempty"`
	Browser    string `bson:"browser,omitempty" json:"browser,omitempty"`
	ScreenSize string `bson:"screen_size,omitempty" json:"screen_size,omitempty"`
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Content   string             `bson:"content" json:"content"`
	Sender    string             `bson:"sender" json:"sender"` // user, assistant, system
	Type      string             `bson:"type" json:"type"`     // text, image, file, quick_reply
	Metadata  MessageMetadata    `bson:"metadata,omitempty" json:"metadata,omitempty"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Read      bool               `bson:"read" json:"read"`
}

type MessageMetadata struct {
	Intent      string   `bson:"intent,omitempty" json:"intent,omitempty"`
	Confidence  float64  `bson:"confidence,omitempty" json:"confidence,omitempty"`
	Entities    []Entity `bson:"entities,omitempty" json:"entities,omitempty"`
	IsTemplate  bool     `bson:"is_template" json:"is_template"`
	TemplateID  string   `bson:"template_id,omitempty" json:"template_id,omitempty"`
	Attachments []string `bson:"attachments,omitempty" json:"attachments,omitempty"`
}

type Entity struct {
	Type  string `bson:"type" json:"type"`
	Value string `bson:"value" json:"value"`
	Start int    `bson:"start" json:"start"`
	End   int    `bson:"end" json:"end"`
}

// Template Schema (for MINILMv2 templates)
type Template struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Type        string             `bson:"type" json:"type"` // greeting, faq, response, escalation
	Category    string             `bson:"category" json:"category"`
	Content     string             `bson:"content" json:"content"`
	Variables   []string           `bson:"variables,omitempty" json:"variables,omitempty"`
	Model       string             `bson:"model" json:"model"` // MINILMv2, etc.
	Version     string             `bson:"version" json:"version"`
	Language    string             `bson:"language" json:"language"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	UsageCount  int                `bson:"usage_count" json:"usage_count"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Analytics Schema
type AnalyticsEvent struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	EventType string                 `bson:"event_type" json:"event_type"`
	WidgetID  primitive.ObjectID     `bson:"widget_id" json:"widget_id"`
	UserID    primitive.ObjectID     `bson:"user_id" json:"user_id"`
	VisitorID string                 `bson:"visitor_id,omitempty" json:"visitor_id,omitempty"`
	PageURL   string                 `bson:"page_url,omitempty" json:"page_url,omitempty"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	Metadata  map[string]interface{} `bson:"metadata" json:"metadata"`
}

// Collection names
const (
	CollectionUsers         = "users"
	CollectionWidgets       = "widgets"
	CollectionConversations = "conversations"
	CollectionTemplates     = "templates"
	CollectionAnalytics     = "analytics_events"
)
