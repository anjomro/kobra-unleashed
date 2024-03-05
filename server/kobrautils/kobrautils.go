package kobrautils

import (
	"fmt"
	"os"
	"time"

	"github.com/anjomro/kobra-unleashed/server/structs"
	"github.com/google/uuid"
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
		// If not, read the file
		file, err := os.ReadFile("/user/ac_mqtt_connect_info")
		if err != nil {
			return "", fmt.Errorf("error reading printer ID")
		}

		// From 0000080 to 0000090
		printerID = string(file[128:160])
		return printerID, nil
	}
}

func GetPrinterModel() (string, error) {
	if printerModel != "" {
		return printerModel, nil
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

// Initialize the MqttPayload struct
func NewMqttPayload(typeStr string, actionStr string, data interface{}) *structs.MqttPayload {
	msgID := uuid.New().String()
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	return &structs.MqttPayload{
		MsgID:     msgID,
		Timestamp: timestamp,
		Type:      typeStr,
		Action:    actionStr,
		Data:      data,
	}
}
