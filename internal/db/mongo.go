package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"guitar_processor/internal/config"
	"time"
)

func NewMongoClient(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := options.Credential{Username: cfg.MongoUser, Password: cfg.MongoPass}

	return mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoUri).SetAuth(auth))
}

func NewMongoDatabase(client *mongo.Client, cfg *config.Config) *mongo.Database {
	return client.Database()
}
