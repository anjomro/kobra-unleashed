package routes

import (
	"bufio"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"github.com/anjomro/kobra-unleashed/server/kobrautils"
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
		fmt.Println("Printing file", filename)
	}

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

func getFilesGET(ctx *fiber.Ctx) error {
	// Return a list of files in /mnt/UDISK and /mnt/exUDISK

	// Get files
	files, err := kobrautils.ListFiles()
	if err != nil {
		slog.Error("ListFiles", "err", err.Error())
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Error getting local files",
		})
	}

	// Return the files as json [{}]
	return ctx.Status(200).JSON(files)

}

func getFileGET(ctx *fiber.Ctx) error {
	// Detect if listLocal or listUdisk
	// Get ?pathType and ?path
	pathType := ctx.Params("pathtype")
	filename := ctx.Params("filename")

	filename, err := url.QueryUnescape(filename)
	if err != nil {
		slog.Error("UrlDecode", "err", err.Error())
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Error decoding filename",
		})
	}

	// If pathType is not listLocal or listUdisk, return 400
	if pathType != "local" && pathType != "usb" {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid pathType",
		})
	}

	// If filename is empty, return 400
	if filename == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	if kobrautils.CheckName(filename) {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	// Only allow .gcode files
	if !strings.HasSuffix(filename, ".gcode") {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	var path string
	if pathType == "local" {
		path = "/mnt/UDISK/"
	} else if pathType == "usb" {
		path = "/mnt/exUDISK/"
	} else {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid pathtype",
		})
	}

	// Make a buffer of 1kb
	buff := make([]byte, 1024)

	// Open the file
	file, err := os.Open(path + filename)
	if err != nil {
		slog.Error("OpenFile", "err", err.Error())
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Error opening file",
		})
	}

	defer file.Close()

	// Set the content type to application/octet-stream
	ctx.Set(fiber.HeaderContentType, fiber.MIMEOctetStream)

	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {

		for {
			// Read the file into the buffer
			n, err := file.Read(buff)
			if err != nil {
				break
			}

			if n == 0 {
				break
			}

			// Write the buffer to the response
			_, err = w.Write(buff[:n])
			if err != nil {
				break
			}

			// Flush the buffer
			err = w.Flush()
			if err != nil {
				break
			}
		}
	})

	return ctx.SendStatus(200)
}

func moveFileGET(ctx *fiber.Ctx) error {
	pathType := ctx.Params("pathtype")
	filename := ctx.Params("filename")
	topath := ctx.Params("topath")

	if pathType != "local" && pathType != "usb" {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid pathtype",
		})
	}

	if topath != "local" && topath != "usb" {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid topath",
		})
	}

	if filename == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	if kobrautils.CheckName(filename) {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	// Url decode filename
	filename, err := kobrautils.UrlDecode(filename)
	if err != nil {
		slog.Error("UrlDecode", "err", err.Error())
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Error decoding filename",
		})
	}

	// Move the file
	// If topath is local, move from usb to local
	// If topath is usb, move from local to usb
	if topath == "local" {
		// Move from usb to local
		err := kobrautils.MoveFile("/mnt/exUDISK/"+filename, "/mnt/UDISK/"+filename)
		if err != nil {
			slog.Error("MoveFile", "err", err.Error())
			return ctx.Status(500).JSON(fiber.Map{
				"error": "Error moving file",
			})
		}
	} else if topath == "usb" {
		// Move from local to usb
		err := kobrautils.MoveFile("/mnt/UDISK/"+filename, "/mnt/exUDISK/"+filename)
		if err != nil {
			slog.Error("MoveFile", "err", err.Error())
			return ctx.Status(500).JSON(fiber.Map{
				"error": "Error moving file",
			})
		}
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "File moved",
	})
}

func deleteFileDELETE(ctx *fiber.Ctx) error {
	pathType := ctx.Params("pathtype")
	filename := ctx.Params("filename")

	if pathType != "local" && pathType != "usb" {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid pathtype",
		})
	}

	if filename == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	if kobrautils.CheckName(filename) {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	// Delete the file
	if pathType == "local" {
		err := kobrautils.DeleteFile("local", filename)
		if err != nil {
			slog.Error("DeleteFile", "err", err.Error())
			return ctx.Status(500).JSON(fiber.Map{
				"error": "Error deleting file",
			})
		}
	} else if pathType == "usb" {
		err := kobrautils.DeleteFile("usb", filename)
		if err != nil {
			slog.Error("DeleteFile", "err", err.Error())
			return ctx.Status(500).JSON(fiber.Map{
				"error": "Error deleting file",
			})
		}
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "File deleted",
	})
}
