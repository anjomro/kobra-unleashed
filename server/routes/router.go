package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Create a router with base path /api
	router := app.Group("/api")

	// /api/
	router.Get("/", indexHandler)
}
