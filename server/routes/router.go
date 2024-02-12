package routes

import (
	"github.com/anjomro/kobra-unleashed/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Create a router with base path /api
	router := app.Group("/api")

	// /api/
	router.Get("/", middleware.AuthHandler, indexHandler)
}
