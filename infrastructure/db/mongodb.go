package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	clientOpts := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("❌ Error connecting to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("❌ Error pinging MongoDB: %w", err)
	}

	fmt.Println("✅ Connected to MongoDB Atlas!")
	return &MongoDB{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}

func (m *MongoDB) Close() {
	if err := m.Client.Disconnect(context.Background()); err != nil {
		log.Println("❌ Error closing MongoDB connection:", err)
	} else {
		fmt.Println("❌ MongoDB connection closed.")
	}
}
