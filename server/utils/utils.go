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

// Make a config method that has Read Config and Save Config

type Config struct {
	Settings structs.Settings
}

func (c *Config) Read() error {
	filename := "/user/settings.json"

	configfile, err := os.Open(filename)
	if err != nil {
		slog.Error("Error opening file", err)
	}

	defer configfile.Close()

	// Read json file
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("Error reading file", err)
		return err
	}

	// Unmarshal json
	err = json.Unmarshal(byteValue, &c.Settings)
	if err != nil {
		slog.Error("Error unmarshalling json", err)
		return err
	}

	return nil
}

func (c *Config) Save() error {
	filename := "/user/settings.json"

	// Write settings to settings.json

	configfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		slog.Error("Error opening file", err)
		return err
	}
	defer configfile.Close()

	// Marshal json
	byteValue, err := json.Marshal(c.Settings)
	if err != nil {
		slog.Error("Error marshalling json", err)
		return err
	}

	// Write to file
	_, err = configfile.Write(byteValue)
	if err != nil {
		slog.Error("Error writing to file", err)
		return err
	}

	return nil
}

func SetConfig() *Config {
	var settings structs.Settings
	c := Config{settings}
	return &c
}
