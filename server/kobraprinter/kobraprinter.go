package kobraprinter

import (
	"fmt"

	"github.com/anjomro/kobra-unleashed/server/kobrautils"
	"github.com/anjomro/kobra-unleashed/server/mqtt"
	"github.com/anjomro/kobra-unleashed/server/structs"
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

//	{
//		"type": "print",
//		"action": "update",
//		"data": {
//		  "taskid": "0",
//		  "settings": {
//			"target_nozzle_temp": 0,
//			"target_hotbed_temp": 0,
//			"fan_speed_pct": 1,
//			"print_speed_mode": 1
//		  }
//		}
//	  }

func UpdatePrintSettings(taskid string, settings structs.PrintSettings) error {
	payld := kobrautils.NewMqttPayload("print", "update", map[string]interface{}{"taskid": taskid, "settings": settings})

	err := mqtt.SendCommand(payld)

	if err != nil {
		return err
	} else {
		return nil
	}
}
