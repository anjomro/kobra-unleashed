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

	// Create a api with base path /api and protect it with the middleware
	api := app.Group("/api").Use(middleware.AuthHandler)

	api.Get("/version", versionHandler)
	api.Get("/user", GetUserInfo)
	api.Post("/print", PrintHandler).Post("/print/:taskid/cancel", CancelPrintHandler).Post("/print/:taskid/pause", PausePrintHandler).Post("/print/:taskid/resume", ResumePrintHandler).Get("/print/query/:taskid", PrintQueryHandler)
	api.Put("/printer/settings", SetPrinterSettingsHandler)
	api.Get("/printer/status", GetPrinterStatusHandler)
	api.Get("/printer/log", GetLogHandler)

	// /api/files/
	filehandler := api.Group("/files")
	filehandler.Post("/local", localFilesHandlerPOST)
	filehandler.Post("/sdcard", sdcardFilesHandlerPOST)

	// /api/files/:pathType/:path
	filehandler.Get("/", getFilesGET)
	filehandler.Get("/:pathtype/:filename", getFileGET).Delete("/:pathtype/:filename", deleteFileDELETE)
	filehandler.Get("/:pathtype/:filename/:topath", moveFileGET)
}
