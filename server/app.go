package main

import (
	"log"
	"log/slog"
	"os"
	"runtime"

	"github.com/anjomro/kobra-unleashed/server/routes"
	"github.com/anjomro/kobra-unleashed/server/utils"
	"github.com/gofiber/contrib/websocket"
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

	utils.CheckSetup()

	appPort := utils.GetEnv("APP_PORT", "3000")

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		slog.Info("Request:", c.Method(), c.Path())
		return c.Next()
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	routes.SetupRoutes(app)

	// Setup a middleware that prints the request method and path and all that is requested

	app.Listen(":" + appPort)
}
