package db

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	once   sync.Once
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}
}

func InitDatabase() {
	// Initialize collections and indexes
	initUsers()
	initPosts()
	log.Println("Database initialized successfully")
}

func initUsers() {
	collection := ConnectUsers()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create unique index for username
	_, err := collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Printf("Warning: Failed to create username index: %v", err)
	}

	// Create unique index for email
	_, err = collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Printf("Warning: Failed to create email index: %v", err)
	}
}

func initPosts() {
	collection := ConnectPosts()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create index for user_id for faster queries
	_, err := collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
	)
	if err != nil {
		log.Printf("Warning: Failed to create user_id index: %v", err)
	}

	// Create index for created_at for sorting
	_, err = collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	)
	if err != nil {
		log.Printf("Warning: Failed to create created_at index: %v", err)
	}
}

// connectDB creates a singleton MongoDB client
func connectDB() *mongo.Client {
	once.Do(func() {
		mongoURI := os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			log.Fatal("MONGODB_URI not set in environment")
		}

		if !strings.HasPrefix(mongoURI, "mongodb://") && !strings.HasPrefix(mongoURI, "mongodb+srv://") {
			mongoURI = "mongodb://" + mongoURI
		}

		clientOptions := options.Client().
			ApplyURI(mongoURI).
			SetConnectTimeout(10 * time.Second).
			SetServerSelectionTimeout(5 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal("Failed to connect to MongoDB:", err)
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal("Failed to ping MongoDB:", err)
		}

		log.Println("Successfully connected to MongoDB")
	})
	return client
}

func ConnectUsers() *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "demo-web-server-2"
	}
	return connectDB().Database(dbName).Collection("users")
}

func ConnectPosts() *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "demo-web-server-2"
	}
	return connectDB().Database(dbName).Collection("posts")
}
