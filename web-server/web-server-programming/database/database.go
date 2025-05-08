package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents the user model
type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

// DBClient represents MongoDB client connection
type DBClient struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewDBClient creates a new MongoDB client
func NewDBClient() (*DBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb+srv://duynghia22302:duynghia123@demo-web-server.e1uwxto.mongodb.net/?retryWrites=true&w=majority&appName=demo-web-server")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return &DBClient{
		client: client,
		db:     client.Database("demo-web-server"),
	}, nil
}

// Close disconnects from MongoDB
func (c *DBClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.Disconnect(ctx)
}

// CreateUser inserts a new user into the database
func (c *DBClient) CreateUser(user User) error {
	collection := c.db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

// GetUser retrieves a user by username
func (c *DBClient) GetUser(username string) (*User, error) {
	collection := c.db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database
func (c *DBClient) GetAllUsers() ([]User, error) {
	collection := c.db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %v", err)
	}
	return users, nil
}
