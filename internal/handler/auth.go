package handler

import (
	"github.com/gofiber/fiber/v2"
	"go-lost-found/internal/models"
	"go-lost-found/internal/repository"
	"go-lost-found/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=user admin"`
	FcmToken string `json:"fcmToken" validate:"required"`
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     false,
			"statusCode": "400",
			"message":    "Invalid request body: " + err.Error(),
			"data":       nil,
		})
	}

	exists, err := repository.CheckUserExists(req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":     false,
				"statusCode": "400",
				"message":    "Error checking user",
				"data":       nil,
			})
	}
	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":     false,
				"statusCode": "400",
				"message":    "Email already exists",
				"data":       nil,
			})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":     false,
				"statusCode": "500",
				"message":    "Error hashing password: " + err.Error(),
				"data":       nil,
			})
	}

	user := models.NewUser(req.Email, string(hashed), req.Role)

	if err := repository.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":     false,
				"statusCode": "500",
				"message":    "Failed to create user: " + err.Error(),
				"data":       nil,
			})
	}

	user.Password = "" // hide password
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     true,
		"statusCode": "201",
		"message":    "User registered successfully",
		"data":       user,
	})
}

func Login(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	exists, err := repository.CheckUserExists(req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":     false,
				"statusCode": "400",
				"message":    "Error checking user",
				"data":       nil,
			})
	}
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":     false,
				"statusCode": "400",
				"message":    "User not found",
				"data":       nil,
			})
	}

	user, err := repository.FindUserByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":     false,
				"statusCode": "400",
				"message":    "Error finding user: " + err.Error(),
				"data":       nil,
			})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":     false,
				"statusCode": "400",
				"message":    "Invalid password: " + err.Error(),
				"data":       nil,
			})
	}
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":     false,
				"statusCode": "500",
				"message":    "Error generating token: " + err.Error(),
				"data":       nil,
			})
	}
	user.BearerToken = token

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "User logged in successfully",
		"data":       user,
	})

}
