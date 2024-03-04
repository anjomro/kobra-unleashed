package kobraprinter

import (
	"fmt"

	"github.com/anjomro/kobra-unleashed/server/kobrautils"
	"github.com/anjomro/kobra-unleashed/server/mqtt"
)

// Map string interface short

func ListFiles(pathType string, path string) error {
	if pathType != "listLocal" && pathType != "listUdisk" {
		return fmt.Errorf("invalid pathType: %s", pathType)
	}

	// If path is empty, set it to root /
	if path == "" {
		path = "/"
	}

	payld := kobrautils.NewMqttPayload("file", pathType, map[string]interface{}{"path": path})

	err := mqtt.SendCommand(payld)

	if err != nil {
		return err
	} else {
		return nil
	}
}
