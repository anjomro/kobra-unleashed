package routes

import (
	"log/slog"
	"os"
	"path/filepath"

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

// /api/files/local
func localFilesHandlerGET(ctx *fiber.Ctx) error {
	// Get all .gcode files in /mnt/UDISK. Enumerate all files in the directory and return the ones that end with .gcode

	// Get all files in /mnt/UDISK
	sdcardDir := "/mnt/UDISK"

	if utils.IsDev() {
		sdcardDir = "dev/sdcard"
	}

	// Get all files in /mnt/UDISK and /mnt/exUDISK
	files := []File{}
	err := filepath.Walk(sdcardDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a .gcode file
		if filepath.Ext(path) == ".gcode" {
			// Get the file info
			fileInfo, err := os.Stat(path)
			if err != nil {
				return err
			}

			// Create a new file struct
			file := File{
				Name:     fileInfo.Name(),
				Path:     path,
				Type:     "gcode",
				TypePath: []string{"gcode"},
				Hash:     "hash",
				Size:     int(fileInfo.Size()),
				Date:     int(fileInfo.ModTime().Unix()),
				Origin:   "local",
				Refs: struct {
					Resource string `json:"resource"`
					Download string `json:"download"`
				}{
					Resource: "resource",
					Download: "download",
				},
				GcodeAnalysis: struct {
					EstimatedPrintTime int `json:"estimatedPrintTime"`
					Filament           struct {
						Length int     `json:"length"`
						Volume float64 `json:"volume"`
					} `json:"filament"`
				}{
					EstimatedPrintTime: 0,
					Filament: struct {
						Length int     `json:"length"`
						Volume float64 `json:"volume"`
					}{
						Length: 0,
						Volume: 0,
					},
				},
				Print: struct {
					Failure int `json:"failure"`
					Success int `json:"success"`
					Last    struct {
						Date    int  `json:"date"`
						Success bool `json:"success"`
					} `json:"last"`
				}{
					Failure: 0,
					Success: 0,
					Last: struct {
						Date    int  `json:"date"`
						Success bool `json:"success"`
					}{
						Date:    0,
						Success: false,
					},
				},
			}

			// Append the file to the files slice
			files = append(files, file)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return ctx.Status(404).JSON(fiber.Map{
			"error": "No files found",
		})
	}

	// Return the files
	return ctx.JSON(files)
}

// /api/files/sdcard
func sdcardFilesHandlerGET(ctx *fiber.Ctx) error {
	// Get all .gcode files in /mnt/exUDISK. Enumerate all files in the directory and return the ones that end with .gcode

	// Get all files in /mnt/exUDISK
	sdcardDir := "/mnt/exUDISK"

	if utils.IsDev() {
		sdcardDir = "dev/sdcard"
	}

	// Get all files in /mnt/exUDISK
	files := []File{}
	err := filepath.Walk(sdcardDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a .gcode file
		if filepath.Ext(path) == ".gcode" {
			// Get the file info
			fileInfo, err := os.Stat(path)
			if err != nil {
				return err
			}

			// Create a new file struct
			file := File{
				Name:     fileInfo.Name(),
				Path:     path,
				Type:     "gcode",
				TypePath: []string{"gcode"},
				Hash:     "hash",
				Size:     int(fileInfo.Size()),
				Date:     int(fileInfo.ModTime().Unix()),
				Origin:   "sdcard",
				Refs: struct {
					Resource string `json:"resource"`
					Download string `json:"download"`
				}{
					Resource: "resource",
					Download: "download",
				},
				GcodeAnalysis: struct {
					EstimatedPrintTime int `json:"estimatedPrintTime"`
					Filament           struct {
						Length int     `json:"length"`
						Volume float64 `json:"volume"`
					} `json:"filament"`
				}{
					EstimatedPrintTime: 0,
					Filament: struct {
						Length int     `json:"length"`
						Volume float64 `json:"volume"`
					}{
						Length: 0,
						Volume: 0,
					},
				},
				Print: struct {
					Failure int `json:"failure"`
					Success int `json:"success"`
					Last    struct {
						Date    int  `json:"date"`
						Success bool `json:"success"`
					} `json:"last"`
				}{
					Failure: 0,
					Success: 0,
					Last: struct {
						Date    int  `json:"date"`
						Success bool `json:"success"`
					}{
						Date:    0,
						Success: false,
					},
				},
			}

			// Append the file to the files slice
			files = append(files, file)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return ctx.Status(404).JSON(fiber.Map{
			"error": "No files found",
		})
	}

	// Return the files
	return ctx.JSON(files)
}

func uploadFileHandler(ctx *fiber.Ctx, savePath string) error {
	// Get multipart form-data
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

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
