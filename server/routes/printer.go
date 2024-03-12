package routes

import (
	"io"
	"os"

	"github.com/anjomro/kobra-unleashed/server/kobraprinter"
	"github.com/anjomro/kobra-unleashed/server/kobrautils"
	"github.com/anjomro/kobra-unleashed/server/mqtt"
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

func GetPrinterStatusHandler(c *fiber.Ctx) error {
	// Post the printer status from mqtt
	payld := kobrautils.NewMqttPayload("status", "query", nil)
	err := mqtt.SendCommand(payld, "status")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	} else {
		return c.SendStatus(200)
	}
}

type LogRequest struct {
	LogType string `json:"logtype"`
}

func returnLogResponse(c *fiber.Ctx, file string) error {
	// Check if the file exists
	if fileinfo, err := os.Stat(file); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{"error": "Log file does not exist"})
	} else {
		// return c.SendFile(file)
		// Read the file and return the content
		file, err := os.Open(file)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error reading log file"})
		}
		defer file.Close()

		// Read the file into a string
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error reading log file"})
		}

		return c.JSON(fiber.Map{"file_size": fileinfo.Size(), "modified_at": fileinfo.ModTime(), "content": string(fileBytes)})
	}
}

func GetLogHandler(c *fiber.Ctx) error {
	// Check what log file to get
	// app or kobra
	var logRequest LogRequest
	err := c.BodyParser(&logRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if logRequest.LogType == "app" {
		// Check if the file exists
		// Return empty file if it does not

		return returnLogResponse(c, "/mnt/UDISK/log")
	}

	if logRequest.LogType == "kobra" {
		// Check if the file exists
		return returnLogResponse(c, "/mnt/UDISK/kobra.log")
	}

	return c.Status(400).JSON(fiber.Map{"error": "Invalid logtype"})
}

func PrintHandler(c *fiber.Ctx) error {
	// Print a file

	// Check if request or fileupload
	if c.FormValue("upload") == "false" {
		// No upload. Just print the file
		// Get the filename
		filename := c.FormValue("file")
		if filename == "" {
			return c.Status(400).JSON(fiber.Map{"error": "No file name provided"})
		}

		if kobrautils.CheckName(filename) {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid filename"})
		}

		shouldCopy := c.FormValue("copy") == "true"
		if shouldCopy {
			// Copy the file from /mnt/exUDISK to /mnt/UDISK
			err := kobrautils.CopyFile("/mnt/exUDISK/"+filename, "/mnt/UDISK/"+filename)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}

		// Print the file
		err := kobraprinter.PrintFile(filename)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		} else {
			return c.SendStatus(200)
		}
	} else {
		// File upload
		// Get the file from the form
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
		}

		// Save the file to /mnt/UDISK

		if kobrautils.CheckName(file.Filename) {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid filename"})
		}

		filename := "/mnt/UDISK/" + file.Filename
		err = c.SaveFile(file, filename)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error saving file"})
		}

		// Print the file
		err = kobraprinter.PrintFile(file.Filename)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		} else {
			return c.SendStatus(200)
		}
	}
}
