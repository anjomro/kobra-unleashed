package kobraprinter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/anjomro/kobra-unleashed/server/mqtt"
	"github.com/anjomro/kobra-unleashed/server/utils"
)

var printerID string
var printerModel string

// Map string interface short
type jsn map[string]interface{}

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
	// Get the printer model from the printer
	// Check if /user/printer_max.cfg exists | Kobra 2 Max
	// Check if /user/printer.cfg | Kobra 2 Pro
	// Check if /user/printer_plus.cfg | Kobra 2 Plus

	// If printerModel is already set, return it

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

func createPayload(data map[string]interface{}) ([]byte, error) {

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %v", err)
	}

	return payloadBytes, nil
}

// SendMQTTCommand sends a command to the printer with json payload
func SendMQTTCommand(payload []byte) error {
	mqttclinet := *mqtt.GetMQTTClient()

	// Check if the MQTT client is connected
	if !mqttclinet.IsConnected() {
		return fmt.Errorf("MQTT client is not connected")
	}

	printerModel, err := GetPrinterModel()
	if err != nil {
		return err
	}

	printerID, err := GetPrinterID()
	if err != nil {
		return err
	}

	// Convert map string interface to interface

	// anycubic/anycubicCloud/v1/server/printer/<PRINTER_MODEL_ID>/<PRINTER_ID>/response
	token := mqttclinet.Publish(fmt.Sprintf("anycubic/anycubicCloud/v1/server/printer/%s/%s/response", printerModel, printerID), 0, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("[SendMQTTCommand]: error sending mqtt command: %v", token.Error())
	} else {
		fmt.Println("MQTT Command Sent")
	}
	if utils.IsDev() {
		fmt.Println("MQTT URL:", fmt.Sprintf("anycubic/anycubicCloud/v1/server/printer/%s/%s/response", printerModel, printerID))
		jsonpayload, _ := json.Marshal(payload)
		fmt.Println("MQTT Payload:", string(jsonpayload))
	}
	return nil
}

func ListFiles(pathType string, path string) error {
	if pathType != "listLocal" && pathType != "listUdisk" {
		return fmt.Errorf("invalid pathType: %s", pathType)
	}

	// If path is empty, set it to root /
	if path == "" {
		path = "/"
	}

	// payload, err := json.Marshal(jsn{
	// 	"type":   "file",
	// 	"action": pathType,
	// 	"data": jsn{
	// 		"path": path,
	// 	},
	// })
	// if err != nil {
	// 	return fmt.Errorf("[ListFiles]: error marshalling payload: %v", err)
	// }

	payload, err := createPayload(jsn{
		"type":   "file",
		"action": pathType,
		"data": jsn{
			"path": path,
		},
	})
	if err != nil {
		return fmt.Errorf("[ListFiles]: error marshalling payload: %v", err)
	}

	err = SendMQTTCommand(payload)

	if err != nil {
		return fmt.Errorf("[ListFiles]: error sending mqtt command: %v", err)
	}
	return nil
}
