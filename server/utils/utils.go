package utils

import "os"

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
