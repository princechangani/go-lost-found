package handler

import (
	"fmt"
	"go-lost-found/internal/models"
	"go-lost-found/internal/repository"
	services "go-lost-found/internal/service"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetItems(c *fiber.Ctx) error {
	log.Println("GetItems")
	items, err := repository.GetItems()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to get items: " + err.Error(),
			"data":       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Items retrieved successfully",
		"data":       items,
	})
}

func GetItemsByID(c *fiber.Ctx) error {
	id := c.Params("id")
	item, err := repository.GetItemsByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to get item: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Item retrieved successfully",
		"data":       item,
	})
}

func CreateItems(c *fiber.Ctx) error {
	var item models.LostFoundItem

	// Parse normal fields manually
	item.Title = c.FormValue("title")
	item.Description = c.FormValue("description")
	item.Location = c.FormValue("location")
	item.CategoryID = c.FormValue("categoryId")
	item.StatusType = c.FormValue("statusType")
	item.Type = c.FormValue("type")
	item.Time = c.FormValue("time")
	item.Email = c.FormValue("email")

	// Convert date (string → time.Time)
	dateStr := c.FormValue("date")
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr) // Adjust format
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  false,
				"message": "Invalid date format",
			})
		}
		item.Date = parsedDate
	}

	// Convert status (string → int)
	statusStr := c.FormValue("status")
	if statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		item.Status = status
	}

	// Handle image file — key must match frontend: formData.append("file", file)
	file, err := c.FormFile("file")
	if err == nil {
		// Save file or upload to cloud
		fileName := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)
		c.SaveFile(file, fileName)
		item.Image = fileName
	}

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	createdItem, err := repository.CreateItems(item)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to create item: " + err.Error(),
		})
	}

	users := []models.User{}
	users, err = repository.GetUsers()
	if err != nil {
		log.Printf("Warning: failed to get users for notification: %v", err)
	} else {
		fcmTokens := []string{}
		for _, user := range users {
			if user.FcmToken != "" {
				fcmTokens = append(fcmTokens, user.FcmToken)
			}
		}
		if len(fcmTokens) > 0 {
			if notifErr := services.SendMultiNotification(fcmTokens, "Lost Found", "New item added", createdItem.Image); notifErr != nil {
				log.Printf("Warning: failed to send notifications: %v", notifErr)
			}
		}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":     true,
		"statusCode": "201",
		"message":    "Item created successfully",
		"data":       createdItem,
	})
}

func UpdateItems(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.LostFoundItem
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     false,
			"statusCode": "400",
			"message":    "Invalid request body: " + err.Error(),
			"data":       nil,
		})
	}
	item, err := repository.UpdateItems(id, item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to update item: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Item updated successfully",
		"data":       item,
	})
}

func DeleteItems(c *fiber.Ctx) error {
	id := c.Params("id")
	err := repository.DeleteItems(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to delete item: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Item deleted successfully",
		"data":       nil,
	})
}
