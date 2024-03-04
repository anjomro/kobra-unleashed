package main

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/anjomro/kobra-unleashed/server/routes"
	"github.com/anjomro/kobra-unleashed/server/sess"
	"github.com/anjomro/kobra-unleashed/server/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	// Crash if not linux
	if runtime.GOOS != "linux" {
		slog.Error("This program is only supported on Linux")
		os.Exit(1)
	}

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", err)
	}

	appPort := utils.GetEnv("APP_PORT", "3000")

	app := fiber.New()

	sess.SetupSessionStore()

	app.Use(func(c *fiber.Ctx) error {
		slog.Info("Request:", c.Method(), c.Path())
		return c.Next()
	})

	routes.SetupRoutes(app)

	app.Static("/", "/www")

	// Catch all other routes and send back 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("Not Found")
	})

	app.Listen(":" + appPort)
}
