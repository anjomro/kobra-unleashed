package routes

import (
	"github.com/anjomro/kobra-unleashed/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Create a router with base path /api
	router := app.Group("/api").Use(middleware.AuthHandler)
	filehandler := router.Group("/files")

	// /api/
	router.Get("/", indexHandler)
	router.Get("/version", versionHandler)

	// /api/files/
	filehandler.Post("/local", localFilesHandlerPOST)
	filehandler.Post("/sdcard", sdcardFilesHandlerPOST)
}
