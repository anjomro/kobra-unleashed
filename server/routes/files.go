package routes

import (
	"log/slog"

	"github.com/anjomro/kobra-unleashed/server/anyprinter"
	"github.com/anjomro/kobra-unleashed/server/utils"
	"github.com/gofiber/fiber/v2"
)

type File struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Type     string   `json:"type"`
	TypePath []string `json:"typePath"`
	Hash     string   `json:"hash"`
	Size     int      `json:"size"`
	Date     int      `json:"date"`
	Origin   string   `json:"origin"`
	Refs     struct {
		Resource string `json:"resource"`
		Download string `json:"download"`
	} `json:"refs"`
	GcodeAnalysis struct {
		EstimatedPrintTime int `json:"estimatedPrintTime"`
		Filament           struct {
			Length int     `json:"length"`
			Volume float64 `json:"volume"`
		} `json:"filament"`
	} `json:"gcodeAnalysis"`
	Print struct {
		Failure int `json:"failure"`
		Success int `json:"success"`
		Last    struct {
			Date    int  `json:"date"`
			Success bool `json:"success"`
		} `json:"last"`
	} `json:"print"`
}

func uploadFileHandler(ctx *fiber.Ctx, savePath string) error {
	// Get multipart form-data
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	shouldPrint := ctx.FormValue("print") == "true"

	// Check if UDISK or exUDISK
	if utils.IsDev() {
		// Check if the file is being uploaded to /mnt/UDISK or /mnt/exUDISK
		if savePath == "/mnt/UDISK/" {
			savePath = "dev/sdcard/"
		}
		if savePath == "/mnt/exUDISK/" {
			savePath = "dev/uploads/"
		}
	}

	filename := savePath + file.Filename

	// Save the file
	err = ctx.SaveFile(file, filename)
	if err != nil {
		slog.Error("Error saving file", err)
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Error saving file",
		})
	}

	// Print the file
	if shouldPrint {
		// Print the file
		anyprinter.Print()
	}

	// Set headers
	ctx.Set("Location", filename)
	return ctx.Status(201).JSON(fiber.Map{
		"message": "File uploaded",
	})
}

func localFilesHandlerPOST(ctx *fiber.Ctx) error {
	return uploadFileHandler(ctx, "/mnt/UDISK/")
}

func sdcardFilesHandlerPOST(ctx *fiber.Ctx) error {
	return uploadFileHandler(ctx, "/mnt/exUDISK/")
}
