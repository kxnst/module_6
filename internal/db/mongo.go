package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"guitar_processor/internal/config"
	"log"
	"time"
)

func NewMongoClient(cfg *config.Config) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := options.Credential{Username: cfg.MongoUser, Password: cfg.MongoPass}

	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoUri).SetAuth(auth))

	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func NewMongoDatabase(client *mongo.Client, cfg *config.Config) *mongo.Database {
	return client.Database(cfg.MongoDatabase)
}
