package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"guitar_processor/internal"
	"guitar_processor/internal/entity"
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
			in := os.Args
			if len(in) < 4 {
				log.Fatalf("Expected 3 arguments, received %d", len(in))
			}

			user := entity.Effect{Slug: in[1], Name: in[2], DSPType: in[3]}

			_, err := db.Collection(repository.EffectsCollectionName).InsertOne(context.Background(), user)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Effect %s successfullt created", user.Name)
		}),
	)

	app.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
