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

	err := mqtt.SendCommand(pathType, "response", jsn{"path": path})

	if err != nil {
		return err
	} else {
		return nil
	}
}
