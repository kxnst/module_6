package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"guitar_processor/internal/entity"
	"guitar_processor/internal/repository"
)

func TestEffectRepository_PutAndGet(t *testing.T) {
	db := GetTestDatabase(t)
	repo := repository.NewEffectRepository(db)

	// Очистка перед тестом
	_ = db.Collection(repository.EffectsCollectionName).Drop(context.Background())

	e := &entity.Effect{
		Slug:    "test-effect",
		Name:    "Test Effect",
		DSPType: "distortion",
	}

	err := repo.Put(e)
	require.NoError(t, err)

	effects, err := repo.GetEffects()
	require.NoError(t, err)
	require.Len(t, effects, 1)

	assert.Equal(t, e.Slug, effects[0].Slug)
	assert.Equal(t, e.Name, effects[0].Name)
	assert.Equal(t, e.DSPType, effects[0].DSPType)
}
