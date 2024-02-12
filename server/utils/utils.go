package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/anjomro/kobra-unleashed/server/structs"
)

func GetEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func GetPrinterID() string {
	return GetEnv("PRINTER_ID", "printer")
}

func IsDev() bool {
	// If APP_ENV is not set, default to "dev"
	return GetEnv("APP_ENV", "dev") == "dev"
}

func CreateAPIKey() string {
	// Open /dev/urandom to get random bytes.
	f, err := os.Open("/dev/urandom")
	if err != nil {
		slog.Error("Error opening /dev/urandom", err)
	}
	defer f.Close()

	// 32 bytes is more than enough for an API key.
	b := make([]byte, 32)
	_, err = f.Read(b)
	if err != nil {
		slog.Error("Error reading /dev/urandom", err)
	}

	// Convert bytes to hex string.
	return fmt.Sprintf("%x", b)
}

func WriteSettings(settings *structs.Settings) {

	filename := "./settings.json"
	if IsDev() {
		filename = "./settings.json"
	} else {
		filename = "/user/settings.json"
	}

	// Create settings.json
	file, err := os.Create(filename)
	if err != nil {
		slog.Error("Error creating settings.json", err)
	}
	defer file.Close()

	// Marshal settings to json
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		slog.Error("Error marshalling settings to json", err)
	}

	// Write settings to settings.json
	_, err = file.Write(settingsJSON)
	if err != nil {
		slog.Error("Error writing settings to settings.json", err)
	}
}
func CheckSetup() {
	// Check if settings.json exists
	filename := "./settings.json"
	if IsDev() {
		filename = "./settings.json"
	} else {
		filename = "/user/settings.json"
	}

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		slog.Info("settings.json does not exist. Creating new settings.json")
		// Create new settings.json with a new API key
		settings := new(structs.Settings)

		// Create new API key
		settings.APIKey = CreateAPIKey()

		// Write settings to settings.json
		WriteSettings(settings)
	}
}
