package routes

import (
	"github.com/anjomro/kobra-unleashed/server/kobraprinter"
	"github.com/anjomro/kobra-unleashed/server/structs"
	"github.com/gofiber/fiber/v2"
)

func SetPrinterSettingsHandler(c *fiber.Ctx) error {
	// Set printer settings

	// Check if any data was sent
	var printerSettings structs.PrintSettingsRequest
	err := c.BodyParser(&printerSettings)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	// Check if taskid is empty
	if printerSettings.Taskid == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Taskid is empty"})
	}

	var settings structs.PrintSettings

	settings.TargetNozzleTemp = printerSettings.TargetNozzleTemp
	settings.TargetHotbedTemp = printerSettings.TargetHotbedTemp
	settings.FanSpeedPct = printerSettings.FanSpeedPct
	settings.PrintSpeedMode = printerSettings.PrintSpeedMode
	settings.Zcompensation = printerSettings.Zcompensation

	// Send the settings to the printer
	err = kobraprinter.UpdatePrintSettings(printerSettings.Taskid, settings)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	} else {
		return c.SendStatus(200)
	}
}
