package repository

import (
	"context"
	"time"

	"go-lost-found/internal/database"
	"go-lost-found/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckUserExists(email string) (bool, error) {
	collection := database.GetCollection("lost_found_item_db", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"email": email}).Err()
	if err == nil {
		return true, nil
	}
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return false, err
}

func GetUsers() ([]models.User, error) {
	collection := database.GetCollection("lost_found_item_db", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var users []models.User
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func CreateUser(user models.User) error {
	collection := database.GetCollection("lost_found_item_db", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	return err
}

func FindUserByEmail(email string) (models.User, error) {
	collection := database.GetCollection("lost_found_item_db", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return user, err
}

func DeleteCategory(id string) error {
	collection := database.GetCollection("lost_found_item_db", "categories")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err

}
