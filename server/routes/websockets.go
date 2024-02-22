package routes

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func websocketHandler(ctx *fiber.Ctx) error {
	// Setup websockets
	if websocket.IsWebSocketUpgrade(ctx) {
		ctx.Locals("allowed", true)
		return ctx.Next()
	}
	return fiber.ErrUpgradeRequired
}

func websocketShellHandler(c *websocket.Conn) {
	// Real time shell /bin/bash

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, msg)
		if err != nil {
			break
		}
	}
}
