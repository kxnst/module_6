package repository

import (
	"context"
	"os"
	"sync"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"guitar_processor/internal/config"
)

var (
	setupOnce sync.Once
	testDB    *mongo.Database
)

func GetTestDatabase(t *testing.T) *mongo.Database {
	t.Helper()

	setupOnce.Do(func() {
		os.Setenv("GO_ENV", "test")
		cfg := config.NewConfig()

		clientOpts := options.Client().ApplyURI(cfg.MongoUri).SetAuth(options.Credential{
			Username:   cfg.MongoUser,
			Password:   cfg.MongoPass,
			AuthSource: "admin",
		})

		client, err := mongo.Connect(context.Background(), clientOpts)
		if err != nil {
			t.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		testDB = client.Database(cfg.MongoDatabase)
	})

	return testDB
}
