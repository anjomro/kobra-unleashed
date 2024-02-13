package main

import (
	"log/slog"

	"github.com/anjomro/kobra-unleashed/server/routes"
	"github.com/anjomro/kobra-unleashed/server/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", err)
	}

	utils.CheckSetup()

	appPort := utils.GetEnv("APP_PORT", "3000")

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		slog.Info("Request:", c.Method(), c.Path())
		return c.Next()
	})

	routes.SetupRoutes(app)

	// Setup a middleware that prints the request method and path and all that is requested

	app.Listen(":" + appPort)
}
