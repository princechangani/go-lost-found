package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go-lost-found/internal/config"
	"go-lost-found/internal/database"
	"go-lost-found/internal/handler"
	"go-lost-found/internal/middleware"
	services "go-lost-found/internal/service"
	"log"
)

func main() {
	// Load environment variables FIRST
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	services.NewEmailService(
		"smtp.gmail.com",
		465,
		"princechangani.dev@gmail.com",
		"mosg tvys putn ltts",
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

	api := app.Group("/api/v1")
	{
		api.Post("user/register", handler.Register)
		api.Post("user/login", handler.Login)

		// Protected routes
		protected := api.Group("/category", middleware.JWTMiddleware())
		{
			protected.Post("/update/:id", handler.UpdateCategory)
			protected.Get("/all", handler.GetCategories)
			protected.Get("/:id", handler.GetCategoryById)
			protected.Post("", handler.CreateCategory)
			protected.Post("/:id", handler.DeleteCategory)
		}

		protected = api.Group("/items", middleware.JWTMiddleware())
		{
			protected.Post("/update/:id", handler.UpdateItems)
			protected.Get("/all", handler.GetItems)
			protected.Get("/:id", handler.GetItemsByID)
			protected.Post("", handler.CreateItems)
			protected.Get("/:id", handler.DeleteItems)
		}
	}
	log.Fatal(app.Listen(":8090"))
}
