package routes

import (
	"github.com/anjomro/kobra-unleashed/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Create a router with base path /api
	router := app.Group("/api")
	filehandler := router.Group("/files")

	// /api/
	router.Get("/", middleware.AuthHandler, indexHandler)
	router.Get("/version", middleware.AuthHandler, versionHandler)

	// /api/files/
	filehandler.Post("/local", middleware.AuthHandler, localFilesHandlerPOST)
	filehandler.Post("/sdcard", middleware.AuthHandler, sdcardFilesHandlerPOST)
}
