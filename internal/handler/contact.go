package handler

import (
	"fmt"
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
		return err
	}

	if err := repository.CreateContact(&contact); err != nil {
		return err
	}

	var item models.LostFoundItem

	item, err := repository.GetItemsByID(contact.ItemID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if item.Type == "Lost" {
		item.Type = "Found"
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
		item.Type = "Claimed"
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

	return c.JSON(contact)
}

func GetContacts(c *fiber.Ctx) error {
	contacts, err := repository.GetContacts()
	if err != nil {
		return err
	}

	return c.JSON(contacts)
}

func GetContactById(c *fiber.Ctx) error {
	id := c.Params("id")
	contact, err := repository.GetContactById(id)
	if err != nil {
		return err
	}

	return c.JSON(contact)
}

func UpdateContact(c *fiber.Ctx) error {
	id := c.Params("id")
	var contact models.Contact
	if err := c.BodyParser(&contact); err != nil {
		return err
	}

	contact.ID = id
	err := repository.CreateContact(&contact)
	if err != nil {
		return err
	}

	return c.JSON(contact)

	return c.JSON(contact)
}

func DeleteContact(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := repository.DeleteContact(id); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Contact deleted successfully",
	})
}
