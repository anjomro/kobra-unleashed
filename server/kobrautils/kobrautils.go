package kobrautils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
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

func MoveFile(src string, dst string) error {
	// Copy the file
	// Don't use os.Rename because it doesn't work across different mount points
	// Copy the file
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	// Remove the file
	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}

func CopyFile(src string, dst string) error {
	// Copy the file
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(pathType, filename string) error {

	filename, err := UrlDecode(filename)
	if err != nil {
		return fmt.Errorf("error decoding filename")
	}

	if strings.HasPrefix(filename, ".") || strings.HasSuffix(filename, ".") || strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") || strings.Contains(filename, "./") || strings.Contains(filename, ".\\") {
		return fmt.Errorf("invalid filename")
	}
	// Delete the file
	if pathType == "local" {
		return os.Remove("/mnt/UDISK/" + filename)
	} else if pathType == "usb" {
		return os.Remove("/mnt/exUDISK/" + filename)
	} else {
		return fmt.Errorf("invalid pathType")
	}

}

type M map[string]interface{}

func MakeJsonWSResp(themap interface{}) ([]byte, error) {
	// Convert the map to json
	jsonResp, err := json.Marshal(themap)
	if err != nil {
		return nil, err
	}

	message := M{
		"message": base64.StdEncoding.EncodeToString(jsonResp),
	}

	jsn, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return jsn, nil
}

func CheckUSB() bool {
	_, err := os.Stat("/dev/sda1")
	return !os.IsNotExist(err)
}

func CheckName(filename string) bool {
	if strings.HasPrefix(filename, ".") || strings.HasSuffix(filename, ".") || strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") || strings.Contains(filename, "./") || strings.Contains(filename, ".\\") {
		return true
	}
	return false
}

type File struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Modified_at time.Time `json:"modified_at"`
	Path        string    `json:"path"`
}

func ListFiles() ([]File, error) {
	// List all the files in the given path

	var fileList []File

	// Only list .gcode files
	localFiles, err := os.ReadDir("/mnt/UDISK")
	if err != nil {
		return nil, err
	}

	// Check if /mnt/exUDISK exists using the errisnotexist function
	if _, err := os.Stat("/mnt/exUDISK"); !os.IsNotExist(err) {
		usbFiles, err := os.ReadDir("/mnt/exUDISK")
		if err != nil {
			return nil, err
		}

		for _, f := range usbFiles {
			fileinfo, err := f.Info()
			if err != nil {
				return nil, err
			}
			if strings.HasSuffix(f.Name(), ".gcode") && fileinfo.Size() > 0 {
				// Get file info using stat

				fileList = append(fileList, File{
					Name:        f.Name(),
					Size:        fileinfo.Size(),
					Modified_at: fileinfo.ModTime(),
					Path:        "usb",
				})
			}
		}
	}

	for _, f := range localFiles {
		fileinfo, err := f.Info()
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(f.Name(), ".gcode") && fileinfo.Size() > 0 {
			// Get file info using stat

			fileList = append(fileList, File{
				Name:        f.Name(),
				Size:        fileinfo.Size(),
				Modified_at: fileinfo.ModTime(),
				Path:        "local",
			})
		}
	}

	return fileList, nil
}

// Remove %20 from the url etc
func UrlDecode(s string) (string, error) {
	// Query escape the string
	return url.QueryUnescape(s)
}

func SetupGCron() {
	// Start the garbage collection cron job
	// Every 10 seconds
	// debug.runGC()
	for {
		time.Sleep(10 * time.Second)
		debug.FreeOSMemory()
	}
}

func GetPrinterStatus() (*structs.PrintState, error) {
	// Get the printer status
	// Open /user/ac_print_state and read the json
	file, err := os.ReadFile("/user/ac_print_state")
	if err != nil {
		return nil, fmt.Errorf("error reading printer state")
	}

	var printState structs.PrintState
	err = json.Unmarshal(file, &printState)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling printer state")
	}

	return &printState, nil
}

func SetPrinterSettings(settings *structs.PrintSettings) error {
	// Set the printer settings
	// Write the json to /user/ac_print_settings
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("error marshalling settings")
	}

	err = os.WriteFile("/user/ac_print_settings", settingsJSON, 0666)
	if err != nil {
		return fmt.Errorf("error writing settings")
	}

	return nil
}
