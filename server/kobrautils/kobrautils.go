package kobrautils

import (
	"fmt"
	"os"

	"github.com/anjomro/kobra-unleashed/server/utils"
)

var printerID string
var printerModel string

func GetPrinterID() (string, error) {
	// Get the printer ID from the printer
	// Open /user/ac_mqtt_connect_info and read first 32 characters

	// If printerID is already set, return it
	if printerID != "" {
		return printerID, nil
	} else {

		if utils.IsDev() {
			// Return 32 character string
			devID := utils.GetEnv("PRINTER_ID", "bfe8c2c0733dcce037b565e0d44d4e23")
			return devID, nil
		} else {
			data, err := os.ReadFile("/user/ac_mqtt_connect_info")
			if err != nil {
				return "", fmt.Errorf("error reading file: %v", err)
			}

			printerID = string(data[:32])
			return printerID, nil
		}
	}
}

func GetPrinterModel() (string, error) {
	if printerModel != "" {
		return printerModel, nil
	} else {
		if utils.IsDev() {
			devModel := utils.GetEnv("PRINTER_MODEL", "20023")
			return devModel, nil
		} else {
			if _, err := os.Stat("/user/printer_max.cfg"); err == nil {
				printerModel = "20023"
			} else if _, err := os.Stat("/user/printer.cfg"); err == nil {
				printerModel = "20021"
			} else if _, err := os.Stat("/user/printer_plus.cfg"); err == nil {
				printerModel = "20022"
			} else {
				return "", fmt.Errorf("error reading printer model")
			}
			return printerModel, nil
		}
	}
}
