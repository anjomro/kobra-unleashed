package routes

import (
	"github.com/anjomro/kobra-unleashed/server/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	ws := app.Group("/ws").Use(middleware.AuthHandler, websocketHandler)
	ws.Get("/shell", websocket.New(websocketShellHandler))
	ws.Get("/info", websocket.New(GetPrinterMessageHandler))

	app.Post("/api/login", LoginHandler).Post("/api/logout", middleware.AuthHandler, LogoutHandler)

	// Create a router with base path /api
	router := app.Group("/api").Use(middleware.AuthHandler)

	router.Get("/version", versionHandler)
	router.Put("/printer/settings", SetPrinterSettingsHandler)

	filehandler := router.Group("/files")
	// /api/files/
	filehandler.Post("/local", localFilesHandlerPOST)
	filehandler.Post("/sdcard", sdcardFilesHandlerPOST)

	// /api/files/:pathType/:path
	filehandler.Get("/", getFilesGET)
}
