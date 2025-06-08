package unit

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"guitar_processor/internal"
	"guitar_processor/internal/config"
	"guitar_processor/internal/effect/service"
	"guitar_processor/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEffectsRegistry_WithMongo(t *testing.T) {
	var db *mongo.Database
	var registry *service.EffectsRegistry
	var cfg *config.Config

	app := fx.New(
		fx.Provide(internal.GetProviders()...),
		fx.Populate(&db, &registry, &cfg),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("Fx start error: %v", err)
	}
	defer app.Stop(context.Background())

	db.Collection("effects").DeleteMany(context.TODO(), nil)

	_, err := db.Collection("effects").InsertOne(context.TODO(), entity.Effect{
		Slug:    "distortion",
		Name:    "FX Distortion",
		DSPType: "distortion",
	})
	assert.NoError(t, err)

	result := registry.GetEffectsInfo()

	var found bool
	for _, eff := range result {
		if eff.Slug == "distortion" {
			found = true
			assert.Equal(t, "FX Distortion", eff.Name)
			assert.Equal(t, "distortion", eff.DspType)
			break
		}
	}
	assert.True(t, found, "distortion effect should be returned")
}
