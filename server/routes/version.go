package routes

import (
	"github.com/gofiber/fiber/v2"
)

type version struct {
	APIVersion    string `json:"api"`
	ServerVersion string `json:"server"`
	Text          string `json:"text"`
}

func versionHandler(ctx *fiber.Ctx) error {
	return ctx.JSON(version{
		APIVersion:    "0.1",
		ServerVersion: "1.3.10",
		Text:          "OctoPrint 1.3.10",
	})
}
