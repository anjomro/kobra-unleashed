package kobraprinter

import (
	"fmt"
	"time"

	"math/rand"

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

	err := mqtt.SendCommand(payld, "file")

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

	err := mqtt.SendCommand(payld, "print")

	if err != nil {
		return err
	} else {
		return nil
	}
}

// str(random.randint(0, 1000000)),

func PrintFile(filename string) error {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	taskid := r.Intn(1000000)

	payld := kobrautils.NewMqttPayload("print", "start", map[string]interface{}{
		"filename":  filename,
		"filepath":  "/",
		"taskid":    fmt.Sprintf("%d", taskid),
		"task_mode": 1,
		"filetype":  1,
	})

	err := mqtt.SendCommand(payld, "print")

	if err != nil {
		return err
	} else {
		return nil
	}
}
