package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const correlationIDHeader = "X-Correlation-ID"

func CorrelationID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		correlationID := c.Get(correlationIDHeader)
		if correlationID == "" {
			correlationID = uuid.NewString()
		}
		c.Locals("correlationID", correlationID)
		c.Set(correlationIDHeader, correlationID)
		return c.Next()
	}
}
