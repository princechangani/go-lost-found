package repository

import (
	"context"
	"go-lost-found/internal/database"
	"go-lost-found/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetItems() ([]models.LostFoundItem, error) {
	collection := database.GetCollection("lost_found_item_db", "lost_item")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var items []models.LostFoundItem

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func GetItemsByID(id string) (models.LostFoundItem, error) {
	collection := database.GetCollection("lost_found_item_db", "lost_item")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var item models.LostFoundItem
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	return item, err

}

func CreateItems(item models.LostFoundItem) (models.LostFoundItem, error) {
	collection := database.GetCollection("lost_found_item_db", "lost_item")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return models.LostFoundItem{}, err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		item.ID = oid.String()
	}

	return item, err

}

func UpdateItems(id string, item models.LostFoundItem) (models.LostFoundItem, error) {
	collection := database.GetCollection("lost_found_item_db", "lost_item")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": item})
	return item, err

}

func DeleteItems(id string) error {
	collection := database.GetCollection("lost_found_item_db", "lost_item")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err

}
