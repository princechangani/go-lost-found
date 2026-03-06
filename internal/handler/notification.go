package handler

import (
	"github.com/gofiber/fiber/v2"
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
