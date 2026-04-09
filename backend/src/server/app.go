package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/martin-aziz/scopra/backend/src/api"
	"github.com/martin-aziz/scopra/backend/src/config"
	"github.com/martin-aziz/scopra/backend/src/middleware"
	"github.com/martin-aziz/scopra/backend/src/services"
)

func New(cfg config.Config, handler *api.Handler, tokenService *services.TokenService) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var fiberError *fiber.Error
			if errors.As(err, &fiberError) {
				return c.Status(fiberError.Code).JSON(fiber.Map{"error": fiberError.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		},
	})

	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOriginURL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))
	app.Use(middleware.CorrelationID())
	app.Use(middleware.Metrics())

	api.RegisterRoutes(app, handler, tokenService)
	return app
}
