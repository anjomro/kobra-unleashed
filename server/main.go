package main

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/anjomro/kobra-unleashed/server/mqtt"
	"github.com/anjomro/kobra-unleashed/server/routes"
	"github.com/anjomro/kobra-unleashed/server/sess"
	"github.com/anjomro/kobra-unleashed/server/utils"
	"github.com/anjomro/kobra-unleashed/server/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	slogfiber "github.com/samber/slog-fiber"
)

func main() {

	// Crash if not linux
	if runtime.GOOS != "linux" {
		slog.Error("This program is only supported on Linux")
		os.Exit(1)
	}

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "errMsg", err.Error())
	}

	// Setup slog to log to file in /mnt/UDISK/kobra.log

	logFile, err := os.OpenFile("/mnt/UDISK/kobra.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Error opening log file", "errMsg", err.Error())
	}
	defer logFile.Close()

	logger := slog.New(slog.NewTextHandler(logFile, nil))

	slog.SetDefault(logger)

	appPort := utils.GetEnv("APP_PORT", "80")

	app := fiber.New()

	app.Use(slogfiber.New(logger))

	sess.SetupSessionStore()

	app.Use(func(c *fiber.Ctx) error {
		slog.Info("Request:", c.Method(), c.Path())
		return c.Next()
	})

	routes.SetupRoutes(app)

	mqtt.SubscribeToPrinter()

	websocket.SetupWebsocket()

	app.Static("/", "/www")

	// Catch all other routes and send back 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("Not Found")
	})

	app.Listen(":" + appPort)

}
