package handler

import (
	"github.com/gofiber/fiber/v2"
	"go-lost-found/internal/models"
	"go-lost-found/internal/repository"
)

func GetCategories(c *fiber.Ctx) error {
	categories, err := repository.GetCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Categories retrieved successfully",
		"data":       categories,
	})
}

func GetCategoryById(c *fiber.Ctx) error {
	id := c.Params("id")
	category, err := repository.GetCategoryById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to get category: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Category retrieved successfully",
		"data":       category,
	})
}

func CreateCategory(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}
	category, err := repository.SaveCategory(category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to create category: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     true,
		"statusCode": "201",
		"message":    "Category created successfully",
		"data":       category,
	})
}

func CreateAllCategory(c *fiber.Ctx) error {
	var categories []models.Category
	if err := c.BodyParser(&categories); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     false,
			"statusCode": "400",
			"message":    "Invalid request body: " + err.Error(),
			"data":       nil,
		})
	}
	categories, err := repository.SaveAllCategory(categories)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to create categories: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     true,
		"statusCode": "201",
		"message":    "Categories created successfully",
		"data":       categories,
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}
	category, err := repository.GetCategoryById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get category",
		})
	}

	category, err = repository.SaveCategory(category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to update category: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Category updated successfully",
		"data":       category,
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repository.DeleteCategory(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     false,
			"statusCode": "500",
			"message":    "Failed to delete category: " + err.Error(),
			"data":       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     true,
		"statusCode": "200",
		"message":    "Category deleted successfully",
		"data":       nil,
	})

}
