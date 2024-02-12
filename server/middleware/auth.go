package middleware

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/anjomro/kobra-unleashed/server/structs"
	"github.com/anjomro/kobra-unleashed/server/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthHandler(ctx *fiber.Ctx) error {
	// Get api key from request header
	apiKey := ctx.Get("X-API-KEY")

	if apiKey == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "No API key provided",
		})
	}

	// Check if api key is valid
	// Check file /user/users.json for valid api keys

	filename := "./settings.json"
	if utils.IsDev() {
		filename = "./settings.json"
	} else {
		filename = "/user/settings.json"
	}

	configfile, err := os.Open(filename)
	if err != nil {
		slog.Error("Error opening file", err)
	}

	defer configfile.Close()

	// Read json file
	byteValue, _ := os.ReadFile(filename)

	// Unmarshal json
	var settings structs.Settings

	// Read json into settings struct
	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		slog.Error("Error unmarshalling json", err)
	}

	// Check if api key is valid
	if apiKey != settings.APIKey {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid API key",
		})
	}

	// Continue to next middleware
	return ctx.Next()
}
