package repository

import (
	"context"
	"testing"

	"guitar_processor/internal/entity"
	"guitar_processor/internal/repository"
)

func TestUserRepository_PutAndGet(t *testing.T) {
	db := GetTestDatabase(t)

	// Очистити колекцію перед тестом
	db.Collection(repository.UserCollectionName).DeleteMany(context.Background(), nil)

	repo := repository.NewUserRepository(db)

	user := &entity.User{
		Login:    "testuser",
		Password: "secret",
	}

	err := repo.Put(user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	foundUser, err := repo.GetUserByLogin("testuser")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if foundUser.Login != "testuser" {
		t.Fatalf("Expected user 'testuser', got %+v", foundUser)
	}
}
