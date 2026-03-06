package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go-lost-found/internal/models"
	"go-lost-found/internal/repository"
	services "go-lost-found/internal/service"
)

type Handler struct {
	Email *services.EmailService
}

func NewHandler(emailService *services.EmailService) *Handler {
	return &Handler{
		Email: emailService,
	}
}
func (h *Handler) CreateContact(c *fiber.Ctx) error {
	var contact models.Contact
	if err := c.BodyParser(&contact); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     false,
			"statusCode": "400",
			"message":    "Invalid request body: " + err.Error(),
			"data":       nil,
		})
	}

	// Validate item exists BEFORE saving contact to avoid orphan records
	item, err := repository.GetItemsByID(contact.ItemID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     false,
			"statusCode": "404",
			"message":    "Item not found: " + err.Error(),
			"data":       nil,
		})
	}

	if err := repository.CreateContact(&contact); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to save contact: " + err.Error(),
			"data":       nil,
		})
	}

	if strings.EqualFold(item.Type, "lost") {
		item.Type = "found"
		body := fmt.Sprintf(`
	<h3>Your item has been found!</h3>
	<p>You can reach the finder at: <b>%s</b></p>`, contact.Email)
		err := h.Email.SendEmail("Your item is found", item.Email, body)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to send email: " + err.Error(),
			})
		}
		var user models.User
		user, err = repository.FindUserByEmail(item.Email)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to find user: " + err.Error(),
			})
		}
		err = services.SendNotification(user.FcmToken, "Your item is found", body, item.Image)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to send notification: " + err.Error(),
			})
		}
	} else {
		item.Type = "claimed"
		body := fmt.Sprintf(`
	<h3>The Item you found has been claimed!</h3>
	<p>You can reach the finder at: <b>%s</b></p>`, contact.Email)
		err := h.Email.SendEmail("The Item you found has been claimed!", item.Email, body)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to send email: " + err.Error(),
			})
		}
		var user models.User
		user, err = repository.FindUserByEmail(item.Email)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to find user: " + err.Error(),
			})
		}
		err = services.SendNotification(user.FcmToken, "The Item you found has been claimed!", body, item.Image)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to send notification: " + err.Error(),
			})
		}

	}

	_, err = repository.UpdateItems(item.ID, item)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to update item: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     true,
		"statusCode": "201",
		"message":    "Contact created successfully",
		"data":       contact,
	})
}

func GetContacts(c *fiber.Ctx) error {
	contacts, err := repository.GetContacts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to get contacts: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Contacts retrieved successfully",
		"data":       contacts,
	})
}

func GetContactById(c *fiber.Ctx) error {
	id := c.Params("id")
	contact, err := repository.GetContactById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to get contact: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Contact retrieved successfully",
		"data":       contact,
	})
}

func UpdateContact(c *fiber.Ctx) error {
	id := c.Params("id")
	var contact models.Contact
	if err := c.BodyParser(&contact); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     false,
			"statusCode": "400",
			"message":    "Invalid request body: " + err.Error(),
			"data":       nil,
		})
	}
	contact.ID = id
	if err := repository.UpdateContact(&contact); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to update contact: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Contact updated successfully",
		"data":       contact,
	})
}

func DeleteContact(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := repository.DeleteContact(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to delete contact: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Contact deleted successfully",
		"data":       nil,
	})
}
