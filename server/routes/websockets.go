package routes

import (
	"bufio"
	"encoding/json"
	"os/exec"

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

type wsMessage struct {
	Message string `json:"message"`
	ErrMess string `json:"errmess,omitempty"`
	ErrCode int    `json:"errcode,omitempty"`
}

type wsCommand struct {
	Command string `json:"command"`
}

func websocketShellHandler(c *websocket.Conn) {
	// Real time shell /bin/bash

	// Process the input and output
	for {
		// Read the message from the browser
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		// Check if ctrl+c was sent and cancel the command
		if string(msg) == "exit\n" {
			break
		}

		// Parse json
		var recCMD wsCommand
		err = json.Unmarshal(msg, &recCMD)
		if err != nil {
			c.WriteJSON(wsMessage{ErrMess: "Invalid JSON", ErrCode: 1})
			continue
		}

		// Check if command is empty
		if recCMD.Command == "" {
			c.WriteJSON(wsMessage{ErrMess: "Command is empty", ErrCode: 2})
			continue
		}

		// Parse the message
		// Execute the command
		cmd := exec.Command("/bin/ash", "-c", recCMD.Command)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			c.WriteJSON(wsMessage{ErrMess: err.Error(), ErrCode: 3})
			continue
		}

		if err := cmd.Start(); err != nil {
			c.WriteJSON(wsMessage{ErrMess: err.Error(), ErrCode: 5})
			continue
		}

		stdo := bufio.NewScanner(stdout)
		for stdo.Scan() {
			// c.WriteMessage(websocket.TextMessage, s.Bytes())
			c.WriteJSON(wsMessage{Message: stdo.Text()})
		}

		if err := cmd.Wait(); err != nil {
			// Check if command not found error
			if err.Error() == "exit status 127" {
				c.WriteJSON(wsMessage{ErrMess: "Command not found", ErrCode: 6})
				continue
			} else {
				c.WriteJSON(wsMessage{ErrMess: err.Error(), ErrCode: 7})
				continue
			}
		}
	}
}
