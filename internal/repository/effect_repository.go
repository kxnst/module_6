package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"guitar_processor/internal/entity"
	"log"
)

const (
	EffectsCollectionName = "effects"
	IndexEffectSlug       = "unique_slug"
)

type EffectRepository struct {
	db *mongo.Database
}

func NewEffectRepository(db *mongo.Database) *EffectRepository {
	return &EffectRepository{db: db}
}

func (er *EffectRepository) Put(effect *entity.Effect) error {
	_, err := er.db.Collection(EffectsCollectionName).InsertOne(context.Background(), effect)

	return err
}

func (er *EffectRepository) GetEffects() ([]*entity.Effect, error) {
	cursor, err := er.db.Collection(EffectsCollectionName).Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Getting effects failed: ")

		return nil, err
	}

	var results []*entity.Effect
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
