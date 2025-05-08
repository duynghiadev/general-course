package db

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	once   sync.Once
)

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
		log.Fatal(err)
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
		log.Fatal(err)
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
		log.Fatal(err)
	}

	// Create index for created_at for sorting
	_, err = collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

// connectDB creates a singleton MongoDB client
func connectDB() *mongo.Client {
	once.Do(func() {
		clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
		var err error
		client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Check the connection
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
	})
	return client
}

func ConnectUsers() *mongo.Collection {
	client := connectDB()
	return client.Database(os.Getenv("DB_NAME")).Collection("users")
}

func ConnectPosts() *mongo.Collection {
	client := connectDB()
	return client.Database(os.Getenv("DB_NAME")).Collection("posts")
}
