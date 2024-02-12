package main

import (
	"log/slog"

	"github.com/anjomro/kobra-unleashed/server/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", err)
	}

	app := fiber.New()

	routes.SetupRoutes(app)

	app.Listen(":3000")
}
