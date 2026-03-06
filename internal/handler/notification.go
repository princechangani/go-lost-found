package handler

import (
	"github.com/gofiber/fiber/v2"
	"go-lost-found/internal/repository"
	services "go-lost-found/internal/service"
)

func SendNotificationHandler(c *fiber.Ctx) error {
	type Req struct {
		Token string `json:"token"`
		Title string `json:"title"`
		Body  string `json:"body"`
		Image string `json:"image"`
	}

	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err := services.SendNotification(req.Token, req.Title, req.Body, req.Image)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Notification sent"})
}

func GetNotificationsHandler(c *fiber.Ctx) error {
	notifications, err := repository.GetNotifications()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to get notifications: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Notifications retrieved successfully",
		"data":       notifications,
	})
}
