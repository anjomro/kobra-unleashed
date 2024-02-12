package utils

import (
	"fmt"
	"log/slog"
	"os"
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
