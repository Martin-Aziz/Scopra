package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/martin-aziz/scopra/backend/src/services"
)

func RequireAccessToken(tokenService *services.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := strings.TrimSpace(c.Get("Authorization"))
		if authorizationHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}

		token := strings.TrimPrefix(authorizationHeader, "Bearer ")
		if token == authorizationHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization scheme"})
		}

		claims, err := tokenService.ParseAndValidate(token, "access")
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("userRole", claims.Role)
		c.Locals("userEmail", claims.Email)
		return c.Next()
	}
}
