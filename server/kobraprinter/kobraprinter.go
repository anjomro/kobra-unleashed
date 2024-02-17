package kobraprinter

import (
	"fmt"

	"github.com/anjomro/kobra-unleashed/server/mqtt"
)

// Map string interface short
type jsn map[string]interface{}

func ListFiles(pathType string, path string) error {
	if pathType != "listLocal" && pathType != "listUdisk" {
		return fmt.Errorf("invalid pathType: %s", pathType)
	}

	// If path is empty, set it to root /
	if path == "" {
		path = "/"
	}

	payload, err := mqtt.CreatePayload(jsn{
		"type":   "file",
		"action": pathType,
		"data": jsn{
			"path": path,
		},
	})
	if err != nil {
		return fmt.Errorf("[ListFiles]: error marshalling payload: %v", err)
	}

	err = mqtt.SendMQTTCommand(payload)

	if err != nil {
		return fmt.Errorf("[ListFiles]: error sending mqtt command: %v", err)
	}
	return nil
}
