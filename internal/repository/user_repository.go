package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"guitar_processor/internal/entity"
)

const (
	UserCollectionName = "users"
	IndexUserLogin     = "unique_login"
)

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Put(user *entity.User) error {
	_, err := ur.db.Collection(UserCollectionName).InsertOne(context.Background(), user)

	return err
}

func (ur *UserRepository) GetUserByLogin(login string) (*entity.User, error) {
	var user *entity.User

	filter := bson.M{"login": login}

	result := ur.db.Collection(UserCollectionName).FindOne(context.Background(), filter)

	if result == nil {
		return nil, fmt.Errorf("user not found")
	}

	err := result.Decode(&user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
