package routes

import (
	"github.com/anjomro/kobra-unleashed/server/middleware"
	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	ws := app.Group("/ws").Use(middleware.AuthHandler, func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	ws.Get("/info", socketio.New(
		func(c *socketio.Websocket) {
			c.Fire("info", make([]byte, 0))
		},
	))

	app.Post("/api/login", LoginHandler).Post("/api/logout", middleware.AuthHandler, LogoutHandler)

	// Create a api with base path /api
	api := app.Group("/api").Use(middleware.AuthHandler)

	api.Get("/version", versionHandler)
	api.Put("/printer/settings", SetPrinterSettingsHandler)
	api.Get("/printer/status", GetPrinterStatusHandler)
	api.Get("/printer/log", GetLogHandler)

	filehandler := api.Group("/files")
	// /api/files/
	filehandler.Post("/local", localFilesHandlerPOST)
	filehandler.Post("/sdcard", sdcardFilesHandlerPOST)

	// /api/files/:pathType/:path
	filehandler.Get("/", getFilesGET)
}
