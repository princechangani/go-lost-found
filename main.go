package main

import (
	"go-lost-found/internal/config"
	"go-lost-found/internal/database"
	"go-lost-found/internal/handler"
	"go-lost-found/internal/middleware"
	services "go-lost-found/internal/service"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables FIRST
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// BUG-012 FIX: Read SMTP credentials from environment variables
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		smtpPort = 465
	}
	emailService := services.NewEmailService(
		os.Getenv("SMTP_HOST"),
		smtpPort,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)

	// Now connect to MongoDB
	database.Connect()

	// Initialize Firebase
	config.InitFirebase()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("MongoDB connection successful! 🚀")
	})
	app.Use(cors.New())

	// Serve uploaded images as static files
	app.Static("/uploads", "./uploads")

	// BUG-014 FIX: Create handler with email service for contact routes
	h := handler.NewHandler(emailService)

	api := app.Group("/api/v1")
	{
		api.Post("/user/register", handler.Register)
		api.Post("/user/login", handler.Login)

		// Category routes — /all must be registered before /:id
		protected := api.Group("/category", middleware.JWTMiddleware())
		{
			protected.Get("/all", handler.GetCategories)
			protected.Get("/:id", handler.GetCategoryById)
			protected.Post("", handler.CreateCategory)
			protected.Put("/:id", handler.UpdateCategory)
			protected.Delete("/:id", handler.DeleteCategory)
		}

		// Item routes — /all must be registered before /:id
		protected = api.Group("/items", middleware.JWTMiddleware())
		{
			protected.Get("/all", handler.GetItems)
			protected.Get("/:id", handler.GetItemsByID)
			protected.Post("", handler.CreateItems)
			protected.Put("/update/:id", handler.UpdateItems)
			protected.Delete("/:id", handler.DeleteItems)
		}

		// Contact routes — /all and /create must be registered before /:id
		protected = api.Group("/contact", middleware.JWTMiddleware())
		{
			protected.Post("/create", h.CreateContact)
			protected.Get("/all", handler.GetContacts)
			protected.Get("/:id", handler.GetContactById)
			protected.Put("/:id", handler.UpdateContact)
			protected.Delete("/:id", handler.DeleteContact)
		}

		// Notification routes
		protected = api.Group("/notifications", middleware.JWTMiddleware())
		{
			protected.Get("/all", handler.GetNotificationsHandler)
			protected.Post("/send", handler.SendNotificationHandler)
		}
	}
	log.Fatal(app.Listen(":8090"))
}
