package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go-lost-found/internal/config"
	"go-lost-found/internal/models"
	"go-lost-found/internal/repository"
	"time"

	"firebase.google.com/go/v4/messaging"
)

func SendNotification(token, title, body, image string) error {
	ctx := context.Background()

	client, err := config.FirebaseApp.Messaging(ctx)
	if err != nil {
		return err
	}

	payload := map[string]string{
		"title":   title,
		"message": body,
		"image":   image,
	}

	data := map[string]string{}
	jsonData, _ := json.Marshal(payload)
	data["message"] = string(jsonData)

	msg := &messaging.Message{
		Token: token,
		Data:  data,
	}

	// Save notification
	notification := models.Notification{
		Title:     title,
		Message:   body,
		Image:     image,
		FcmToken:  token,
		Type:      "Lost",
		Date:      time.Now().Format("2006-01-02"),
		Time:      time.Now().Format("15:04:05"),
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repository.SaveNotification(notification)

	// Send
	response, err := client.Send(ctx, msg)
	if err != nil {
		return err
	}

	fmt.Println("Successfully sent:", response)
	return nil
}
func SendMultiNotification(tokens []string, title, body, image string) error {
	ctx := context.Background()

	client, err := config.FirebaseApp.Messaging(ctx)
	if err != nil {
		return err
	}

	payload := map[string]string{
		"title":   title,
		"message": body,
		"image":   image,
	}

	message := &messaging.MulticastMessage{
		Data:   payload,
		Tokens: tokens,
	}

	response, err := client.SendMulticast(ctx, message)
	if err != nil {
		return err
	}

	// Save notifications
	var notifs []models.Notification
	for _, token := range tokens {
		notifs = append(notifs, models.Notification{
			Title:     title,
			Message:   body,
			Image:     image,
			FcmToken:  token,
			Type:      "Lost",
			Date:      time.Now().Format("2006-01-02"),
			Time:      time.Now().Format("15:04:05"),
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	repository.SaveNotifications(notifs)

	fmt.Println("Success:", response)
	return nil
}
