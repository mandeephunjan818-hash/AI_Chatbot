// mongodb_tester.go - Comprehensive MongoDB testing utilities
package mongoDB_tester

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/gofiles/mongoDB"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SchemaTester tests database schemas and collections
type SchemaTester struct {
	Connector *mongoDB.Connector
	TestDB    *mongo.Database
}

// NewSchemaTester creates a new schema tester
func NewSchemaTester(connector *mongoDB.Connector) (*SchemaTester, error) {
	db, err := connector.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %v", err)
	}
	return &SchemaTester{
		Connector: connector,
		TestDB:    db,
	}, nil
}

// TestConnection performs comprehensive connection tests
func (st *SchemaTester) TestConnection() error {
	log.Println("🧪 Testing MongoDB Connection...")

	// Test basic connectivity
	start := time.Now()
	status := st.Connector.GetStatus()
	elapsed := time.Since(start)

	log.Printf("✅ Connection Status: %+v", status)
	log.Printf("⏱️  Connection latency: %v", elapsed)

	// Test ping
	start = time.Now()
	health := st.Connector.HealthCheck()
	elapsed = time.Since(start)

	log.Printf("🏥 Health Check: %s - %s", health.Status, health.Message)
	log.Printf("⏱️  Ping latency: %v", elapsed)

	if health.Status != "healthy" {
		return fmt.Errorf("connection health check failed: %s", health.Message)
	}

	log.Println("✅ Connection tests completed successfully")
	return nil
}

// TestCollections checks if all required collections exist
func (st *SchemaTester) TestCollections() error {
	log.Println("📁 Testing Database Collections...")

	collections := []string{
		mongoDB.CollectionUsers,
		mongoDB.CollectionWidgets,
		mongoDB.CollectionConversations,
		mongoDB.CollectionTemplates,
		mongoDB.CollectionAnalytics,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collectionNames, err := st.TestDB.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %v", err)
	}

	collectionMap := make(map[string]bool)
	for _, name := range collectionNames {
		collectionMap[name] = true
	}

	missingCollections := []string{}
	for _, required := range collections {
		if !collectionMap[required] {
			missingCollections = append(missingCollections, required)
			log.Printf("⚠️  Missing collection: %s", required)
		} else {
			log.Printf("✅ Collection exists: %s", required)
		}
	}

	if len(missingCollections) > 0 {
		log.Println("📋 Creating missing collections...")
		for _, collection := range missingCollections {
			err := st.TestDB.CreateCollection(ctx, collection)
			if err != nil {
				log.Printf("❌ Failed to create collection %s: %v", collection, err)
			} else {
				log.Printf("✅ Created collection: %s", collection)
			}
		}
	}

	log.Println("✅ Collection tests completed")
	return nil
}

// TestIndexes creates and verifies indexes for each collection
func (st *SchemaTester) TestIndexes() error {
	log.Println("📊 Testing Database Indexes...")

	// Users collection indexes
	userIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"created_at": -1},
		},
		{
			Keys: bson.M{"is_active": 1},
		},
	}

	err := st.createIndexes(mongoDB.CollectionUsers, userIndexes)
	if err != nil {
		return fmt.Errorf("failed to create user indexes: %v", err)
	}

	// Widgets collection indexes
	widgetIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_published", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "type", Value: 1}},
		},
	}

	err = st.createIndexes(mongoDB.CollectionWidgets, widgetIndexes)
	if err != nil {
		return fmt.Errorf("failed to create widget indexes: %v", err)
	}

	// Conversations collection indexes
	conversationIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "widget_id", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "visitor_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "updated_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "assigned_to", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "widget_id", Value: 1},
				{Key: "status", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
	}

	err = st.createIndexes(mongoDB.CollectionConversations, conversationIndexes)
	if err != nil {
		return fmt.Errorf("failed to create conversation indexes: %v", err)
	}

	// Templates collection indexes
	templateIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"type": 1},
		},
		{
			Keys: bson.M{"category": 1},
		},
		{
			Keys: bson.M{"is_active": 1},
		},
		{
			Keys: bson.M{"created_at": -1},
		},
	}

	err = st.createIndexes(mongoDB.CollectionTemplates, templateIndexes)
	if err != nil {
		return fmt.Errorf("failed to create template indexes: %v", err)
	}

	log.Println("✅ Index tests completed successfully")
	return nil
}

func (st *SchemaTester) createIndexes(collectionName string, indexes []mongo.IndexModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := st.TestDB.Collection(collectionName)
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes for %s: %v", collectionName, err)
	}

	log.Printf("✅ Created indexes for %s collection", collectionName)
	return nil
}

// TestCRUDOperations tests CRUD operations on all schema types
func (st *SchemaTester) TestCRUDOperations() error {
	log.Println("🔧 Testing CRUD Operations...")

	// Test Users CRUD
	err := st.testUserCRUD()
	if err != nil {
		return fmt.Errorf("user CRUD test failed: %v", err)
	}

	// Test Widgets CRUD
	err = st.testWidgetCRUD()
	if err != nil {
		return fmt.Errorf("widget CRUD test failed: %v", err)
	}

	// Test Conversations CRUD
	err = st.testConversationCRUD()
	if err != nil {
		return fmt.Errorf("conversation CRUD test failed: %v", err)
	}

	// Test Templates CRUD
	err = st.testTemplateCRUD()
	if err != nil {
		return fmt.Errorf("template CRUD test failed: %v", err)
	}

	log.Println("✅ All CRUD tests completed successfully")
	return nil
}

func (st *SchemaTester) testUserCRUD() error {
	log.Println("  👤 Testing User CRUD operations...")

	collection := st.TestDB.Collection(mongoDB.CollectionUsers)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Clean up any existing test data
	collection.DeleteMany(ctx, bson.M{"email": "test@example.com"})

	// Create (Insert)
	user := mongoDB.User{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashed_password_123",
		FullName:     "Test User",
		Plan:         "pro",
		Settings: mongoDB.UserSettings{
			Theme:              "dark",
			EmailNotifications: true,
			MaxWidgets:         10,
			StorageLimit:       1024,
		},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsActive:   true,
		IsVerified: true,
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}

	userID := result.InsertedID.(primitive.ObjectID)
	log.Printf("    ✅ User inserted with ID: %s", userID.Hex())

	// Read (Find)
	var foundUser mongoDB.User
	err = collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&foundUser)
	if err != nil {
		return fmt.Errorf("find failed: %v", err)
	}

	if foundUser.Email != user.Email {
		return fmt.Errorf("email mismatch: expected %s, got %s", user.Email, foundUser.Email)
	}
	log.Println("    ✅ User found successfully")

	// Update
	update := bson.M{
		"$set": bson.M{
			"full_name":  "Updated Test User",
			"plan":       "enterprise",
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}
	log.Println("    ✅ User updated successfully")

	// Verify update
	err = collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&foundUser)
	if err != nil {
		return fmt.Errorf("verify update failed: %v", err)
	}

	if foundUser.FullName != "Updated Test User" {
		return fmt.Errorf("update verification failed")
	}

	// Delete
	_, err = collection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}
	log.Println("    ✅ User deleted successfully")

	// Verify delete
	count, err := collection.CountDocuments(ctx, bson.M{"_id": userID})
	if err != nil {
		return fmt.Errorf("count failed: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("user still exists after delete")
	}

	log.Println("    ✅ User CRUD operations completed successfully")
	return nil
}

func (st *SchemaTester) testWidgetCRUD() error {
	log.Println("  🎯 Testing Widget CRUD operations...")

	collection := st.TestDB.Collection(mongoDB.CollectionWidgets)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a test user first
	userID := primitive.NewObjectID()

	// Create widget
	widget := mongoDB.Widget{
		UserID:      userID,
		Name:        "Test Widget",
		Type:        "chat",
		Description: "A test widget for CRUD operations",
		Config: mongoDB.WidgetConfig{
			WelcomeMessage:   "Hello! How can I help you?",
			ResponseDelay:    1000,
			MaxMessageLength: 1000,
			Language:         "en",
			Timezone:         "UTC",
			Model:            "MINILMv2",
			Temperature:      0.7,
			MaxTokens:        150,
		},
		Appearance: mongoDB.AppearanceConfig{
			PrimaryColor:   "#3B82F6",
			SecondaryColor: "#1E40AF",
			TextColor:      "#FFFFFF",
			ButtonText:     "Chat with us",
			Icon:           "💬",
			Position:       "bottom-right",
		},
		Behavior: mongoDB.BehaviorConfig{
			AutoOpen:      false,
			DelaySeconds:  5,
			ShowOnMobile:  true,
			ShowOnTablet:  true,
			ShowOnDesktop: true,
		},
		Domains:     []string{"example.com"},
		Placement:   "bottom-right",
		IsActive:    true,
		IsPublished: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := collection.InsertOne(ctx, widget)
	if err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}

	widgetID := result.InsertedID.(primitive.ObjectID)
	log.Printf("    ✅ Widget inserted with ID: %s", widgetID.Hex())

	// Test find by user
	filter := bson.M{"user_id": userID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("count by user failed: %v", err)
	}

	if count != 1 {
		return fmt.Errorf("expected 1 widget for user, got %d", count)
	}

	// Test update
	update := bson.M{
		"$set": bson.M{
			"name":       "Updated Widget Name",
			"is_active":  false,
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": widgetID}, update)
	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	// Clean up
	_, err = collection.DeleteOne(ctx, bson.M{"_id": widgetID})
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}

	log.Println("    ✅ Widget CRUD operations completed successfully")
	return nil
}

func (st *SchemaTester) testConversationCRUD() error {
	log.Println("  💬 Testing Conversation CRUD operations...")

	collection := st.TestDB.Collection(mongoDB.CollectionConversations)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create test IDs
	userID := primitive.NewObjectID()
	widgetID := primitive.NewObjectID()

	// Create conversation
	conversation := mongoDB.Conversation{
		WidgetID:  widgetID,
		UserID:    userID,
		VisitorID: "visitor_123",
		PageURL:   "https://example.com/page",
		Status:    "active",
		Messages: []mongoDB.Message{
			{
				Content:   "Hello, I need help with my order",
				Sender:    "user",
				Type:      "text",
				Timestamp: time.Now().Add(-5 * time.Minute),
				Read:      true,
			},
			{
				Content:   "I'd be happy to help! Can you provide your order number?",
				Sender:    "assistant",
				Type:      "text",
				Timestamp: time.Now().Add(-4 * time.Minute),
				Read:      true,
			},
		},
		CreatedAt: time.Now().Add(-5 * time.Minute),
		UpdatedAt: time.Now(),
	}

	result, err := collection.InsertOne(ctx, conversation)
	if err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}

	convID := result.InsertedID.(primitive.ObjectID)
	log.Printf("    ✅ Conversation inserted with ID: %s", convID.Hex())

	// Test find by widget
	filter := bson.M{"widget_id": widgetID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("count by widget failed: %v", err)
	}

	if count != 1 {
		return fmt.Errorf("expected 1 conversation for widget, got %d", count)
	}

	// Test update with message append
	newMessage := mongoDB.Message{
		Content:   "My order number is ORD-12345",
		Sender:    "user",
		Type:      "text",
		Timestamp: time.Now(),
		Read:      false,
	}

	update := bson.M{
		"$push": bson.M{"messages": newMessage},
		"$set": bson.M{
			"updated_at": time.Now(),
			"status":     "active",
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": convID}, update)
	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	// Clean up
	_, err = collection.DeleteOne(ctx, bson.M{"_id": convID})
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}

	log.Println("    ✅ Conversation CRUD operations completed successfully")
	return nil
}

func (st *SchemaTester) testTemplateCRUD() error {
	log.Println("  📝 Testing Template CRUD operations...")

	collection := st.TestDB.Collection(mongoDB.CollectionTemplates)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create template
	template := mongoDB.Template{
		Name:        "Welcome Greeting",
		Description: "Template for welcoming new visitors",
		Type:        "greeting",
		Category:    "customer_service",
		Content:     "Hello! Welcome to {{company_name}}. How can I assist you today?",
		Variables:   []string{"company_name"},
		Model:       "MINILMv2",
		Version:     "1.0",
		Language:    "en",
		IsActive:    true,
		UsageCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := collection.InsertOne(ctx, template)
	if err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}

	templateID := result.InsertedID.(primitive.ObjectID)
	log.Printf("    ✅ Template inserted with ID: %s", templateID.Hex())

	// Test find by type
	filter := bson.M{"type": "greeting", "is_active": true}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("count by type failed: %v", err)
	}

	if count < 1 {
		return fmt.Errorf("expected at least 1 greeting template, got %d", count)
	}

	// Clean up
	_, err = collection.DeleteOne(ctx, bson.M{"_id": templateID})
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}

	log.Println("    ✅ Template CRUD operations completed successfully")
	return nil
}

// RunAllTests runs all tests in sequence
func (st *SchemaTester) RunAllTests() error {
	log.Println("🚀 Starting Comprehensive MongoDB Tests")
	log.Println("========================================")

	tests := []struct {
		name string
		test func() error
	}{
		{"Connection Test", st.TestConnection},
		{"Collection Test", st.TestCollections},
		{"Index Test", st.TestIndexes},
		{"CRUD Test", st.TestCRUDOperations},
	}

	for _, t := range tests {
		log.Printf("\n📋 Running: %s", t.name)
		log.Println("----------------------------------------")

		err := t.test()
		if err != nil {
			return fmt.Errorf("❌ %s failed: %v", t.name, err)
		}

		log.Printf("✅ %s passed", t.name)
		time.Sleep(500 * time.Millisecond) // Small delay between tests
	}

	log.Println("\n========================================")
	log.Println("🎉 All tests completed successfully!")
	return nil
}
