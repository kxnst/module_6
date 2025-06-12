package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"guitar_processor/internal"
	"guitar_processor/internal/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := fx.New(
		fx.Provide(internal.GetProviders()...),
		fx.Invoke(func(db *mongo.Database) {
			collections := []string{repository.UserCollectionName, repository.EffectsCollectionName}
			for _, collectionName := range collections {
				ctx := context.Background()
				err := db.CreateCollection(ctx, collectionName)
				if err != nil {
					log.Fatalf("Failed to create collection: %v", err)
				}
				fmt.Printf("Collection created: %s\n", collectionName)
			}

			_, err := db.Collection(repository.EffectsCollectionName).Indexes().CreateOne(context.Background(), mongo.IndexModel{
				Keys:    map[string]interface{}{"slug": 1},
				Options: options.Index().SetUnique(true).SetName(repository.IndexEffectSlug),
			})
			if err != nil {
				log.Fatalf("Failed to create index: %v", err)
			}

			_, err = db.Collection(repository.UserCollectionName).Indexes().CreateOne(context.Background(), mongo.IndexModel{
				Keys:    map[string]interface{}{"login": 1},
				Options: options.Index().SetUnique(true).SetName(repository.IndexUserLogin),
			})
			if err != nil {
				log.Fatalf("Failed to create index: %v", err)
			}

			log.Println("Unique indexes created")
		}),
	)

	app.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
