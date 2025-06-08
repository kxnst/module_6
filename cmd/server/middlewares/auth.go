package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"guitar_processor/cmd/server/utils"
)

type AuthMiddleware struct {
	as utils.AuthService
}

func (am *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
	}

	user, err := am.as.ValidateTokenAndGetUser(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	c.Locals("user", user)
	return c.Next()
}

func NewAuthMiddleware(as utils.AuthService) *AuthMiddleware {
	return &AuthMiddleware{as: as}
}
