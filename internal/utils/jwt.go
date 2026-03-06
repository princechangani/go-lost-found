package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret reads the secret lazily so it's always read after godotenv.Load()
func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GenerateToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
}
