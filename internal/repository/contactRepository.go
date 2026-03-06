package repository

import (
	"context"
	"go-lost-found/internal/database"
	"go-lost-found/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func CreateContact(contact *models.Contact) error {
	collection := database.GetCollection("lost_found_item_db", "contacts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, contact)
	return err
}

func GetContacts() ([]models.Contact, error) {
	collection := database.GetCollection("lost_found_item_db", "contacts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var contacts []models.Contact

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &contacts); err != nil {
		return nil, err
	}

	return contacts, nil
}

func GetContactById(id string) (models.Contact, error) {
	collection := database.GetCollection("lost_found_item_db", "contacts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var contact models.Contact

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&contact)
	if err != nil {
		return models.Contact{}, err
	}

	return contact, nil
}

func DeleteContact(id string) error {
	collection := database.GetCollection("lost_found_item_db", "contacts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})

	return err

}
