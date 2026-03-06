package middleware

import (
	"go-lost-found/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// UserIDKey is the key used to store user ID in the Fiber context
const UserIDKey = "userID"

// JWTMiddleware verifies the JWT token and sets the user ID in the context
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing token",
			})
		}

		token, err := utils.ParseToken(tokenStr)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		userID, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		c.Locals(UserIDKey, userID)
		return c.Next()
	}
}

func GetUserIDFromContext(c *fiber.Ctx) (string, bool) {
	userID, ok := c.Locals(UserIDKey).(string)
	return userID, ok
}
