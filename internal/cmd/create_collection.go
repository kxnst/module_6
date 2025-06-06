package main

import (
	"context"
	"fmt"
	"guitar_processor/internal"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"

	"guitar_processor/internal/config"
	"guitar_processor/internal/db"
)

type Params struct {
	fx.In
	DB *mongo.Database
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("⚠️ Usage: create_collection <collection_name>")
	}
	collectionName := os.Args[1]

	app := fx.New(
		fx.Provide(internal.GetProviders()),
		fx.Invoke(func(p Params) {
			ctx := context.Background()
			err := p.DB.CreateCollection(ctx, collectionName)
			if err != nil {
				log.Fatalf("❌ Failed to create collection: %v", err)
			}
			fmt.Printf("✅ Collection created: %s\n", collectionName)
		}),
	)

	app.Run()
}
