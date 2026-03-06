package repository

import (
	"context"
	"go-lost-found/internal/database"
	"go-lost-found/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func SaveNotification(notification models.Notification) error {
	collection := database.GetCollection("lost_found_item_db", "notification")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, notification)
	return err

}

func SaveNotifications(notifications []models.Notification) error {
	collection := database.GetCollection("lost_found_item_db", "notification")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	docs := make([]interface{}, len(notifications))
	for i, v := range notifications {
		docs[i] = v

	}

	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil

}

func GetNotificationByID(id string) error {
	collection := database.GetCollection("lost_found_item_db", "notification")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var notification models.Notification
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&notification)
	return err

}

func GetNotifications() ([]models.Notification, error) {
	collection := database.GetCollection("lost_found_item_db", "notification")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var notification []models.Notification

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &notification); err != nil {
		return nil, err
	}

	return notification, nil
}
