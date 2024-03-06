package routes

import (
	"github.com/anjomro/kobra-unleashed/server/middleware"
	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	ws := app.Group("/ws").Use(func(c *fiber.Ctx) error {
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

	// Create a router with base path /api
	router := app.Group("/api").Use(middleware.AuthHandler)

	router.Get("/version", versionHandler)
	router.Put("/printer/settings", SetPrinterSettingsHandler)
	router.Get("/printer/status", GetPrinterStatusHandler)
	router.Get("/printer/log", GetLogHandler)

	filehandler := router.Group("/files")
	// /api/files/
	filehandler.Post("/local", localFilesHandlerPOST)
	filehandler.Post("/sdcard", sdcardFilesHandlerPOST)

	// /api/files/:pathType/:path
	filehandler.Get("/", getFilesGET)
}
