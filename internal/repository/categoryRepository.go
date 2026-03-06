package repository

import (
	"context"
	"go-lost-found/internal/database"
	"go-lost-found/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCategories() ([]models.Category, error) {
	collection := database.GetCollection("lost_found_item_db", "category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var categories []models.Category

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var category models.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func GetCategoryById(id string) (models.Category, error) {
	collection := database.GetCollection("lost_found_item_db", "category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Category{}, err
	}

	var category models.Category
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&category)
	return category, err

}

func SaveCategory(category models.Category) (models.Category, error) {
	collection := database.GetCollection("lost_found_item_db", "category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, category)
	if err != nil {
		return models.Category{}, err
	}

	// Ensure ID is correctly assigned
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		category.ID = oid.Hex()
	}

	return category, nil
}

func UpdateCategory(id string, category models.Category) (models.Category, error) {
	collection := database.GetCollection("lost_found_item_db", "category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Category{}, err
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": category})
	return category, err
}

func SaveAllCategory(categories []models.Category) ([]models.Category, error) {
	collection := database.GetCollection("lost_found_item_db", "category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert []models.Category → []interface{}
	docs := make([]interface{}, len(categories))
	for i, v := range categories {
		docs[i] = v
	}

	// Insert all
	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	// Assign generated IDs
	for i, id := range result.InsertedIDs {
		if oid, ok := id.(primitive.ObjectID); ok && i < len(categories) {
			categories[i].ID = oid.Hex()
		}
	}

	return categories, nil
}
