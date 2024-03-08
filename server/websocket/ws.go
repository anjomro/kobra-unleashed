package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/anjomro/kobra-unleashed/server/kobraprinter"
	"github.com/anjomro/kobra-unleashed/server/kobrautils"
	"github.com/anjomro/kobra-unleashed/server/mqtt"
	"github.com/gofiber/contrib/socketio"
)

type M map[string]interface{}

const (
	pongWait = 30 * time.Second
)

// map of clients uuid
var clients = make(map[string]*socketio.Websocket)

func SetupWebsocket() {

	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		// Add client id to clinets
		clients[ep.Kws.UUID] = ep.Kws
		slog.Info("Client connected", "clientID", ep.Kws.UUID, "ip", ep.Kws.Conn.RemoteAddr())

		ep.Kws.Conn.SetReadDeadline(time.Now().Add(pongWait))
		ep.Kws.Conn.SetPongHandler(func(string) error {
			ep.Kws.Conn.SetReadDeadline(time.Now().Add(pongWait))
			fmt.Println("Pong received")
			return nil
		})
	})

	socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
		// Remove client id from clients
		delete(clients, ep.Kws.UUID)
	})

	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		ep.Kws.Close()
		// Remove client id from clients
		delete(clients, ep.Kws.UUID)
	})

	type jsonError struct {
		Error string `json:"error"`
	}

	type usbQuery struct {
		IsConnected bool `json:"usb_connected"`
	}

	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {
		// If ping received, send pong
		if string(ep.Data) == "ping" {
			ep.Kws.Emit([]byte("pong"), socketio.TextMessage)
			ep.Kws.Conn.SetReadDeadline(time.Now().Add(pongWait))
		}

		// Try to json decode the message

		var msg M
		if err := json.Unmarshal(ep.Data, &msg); err != nil {
			slog.Error(err.Error())
			return
		}

		// If the message is a move command
		if msg["action"] == "moveToUdisk" {
			filename := msg["filename"].(string)

			// Move file from /mnt/UDISK to /mnt/exUDISK
			if err := kobrautils.MoveFile("/mnt/UDISK/"+filename, "/mnt/exUDISK/"+filename); err != nil {
				slog.Error(err.Error())
				return
			}

			// Trigger a file list update
			err := kobraprinter.ListFiles("listUdisk", "/")
			if err != nil {
				slog.Error(err.Error())
			}

			err = kobraprinter.ListFiles("listLocal", "/")
			if err != nil {
				slog.Error(err.Error())
			}
		}

		if msg["action"] == "moveToLocal" {
			filename := msg["filename"].(string)

			if strings.HasPrefix(filename, ".") || strings.HasSuffix(filename, ".") || strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") || strings.Contains(filename, "./") || strings.Contains(filename, ".\\") {
				slog.Error("Invalid filename")
				// Create json payload
				payload := jsonError{
					Error: "Invalid filename",
				}
				jsonBytes, err := json.Marshal(payload)
				if err != nil {
					slog.Error(err.Error())
					return
				}
				// Send json payload
				ep.Kws.Emit(jsonBytes, socketio.TextMessage)
				return
			}

			// Move file from /mnt/exUDISK to /mnt/UDISK
			if err := kobrautils.MoveFile("/mnt/exUDISK/"+filename, "/mnt/UDISK/"+filename); err != nil {
				slog.Error(err.Error())
				return
			}

			// Trigger a file list update
			err := kobraprinter.ListFiles("listUdisk", "/")
			if err != nil {
				slog.Error(err.Error())
			}

			err = kobraprinter.ListFiles("listLocal", "/")
			if err != nil {
				slog.Error(err.Error())
			}
		}

		if msg["action"] == "deleteFile" {
			filename := msg["filename"].(string)
			filelocation := msg["filelocation"].(string)

			if filelocation == "listLocal" {
				if err := kobrautils.DeleteFile("local", filename); err != nil {
					slog.Error(err.Error())
					return
				}
			} else if filelocation == "listUdisk" {
				if err := kobrautils.DeleteFile("udisk", filename); err != nil {
					slog.Error(err.Error())
					return
				}
			}

			// Trigger a file list update
			err := kobraprinter.ListFiles("listUdisk", "/")
			if err != nil {
				slog.Error(err.Error())
			}

			err = kobraprinter.ListFiles("listLocal", "/")
			if err != nil {
				slog.Error(err.Error())
			}
		}

		if msg["action"] == "check-usb" {
			if kobrautils.CheckUSB() {
				usbQ := usbQuery{
					IsConnected: true,
				}
				resp, err := kobrautils.MakeJsonWSResp(usbQ)
				if err != nil {
					slog.Error(err.Error())
					return
				}
				// Send json payload
				ep.Kws.Emit(resp, socketio.TextMessage)
			} else {
				usbQ := usbQuery{
					IsConnected: false,
				}
				resp, err := kobrautils.MakeJsonWSResp(usbQ)
				if err != nil {
					slog.Error(err.Error())
					return
				}
				// Send json payload
				ep.Kws.Emit(resp, socketio.TextMessage)

				// ep.Kws.Emit(jsonBytes, socketio.TextMessage)
			}
		}
	})

	socketio.On("info", func(ep *socketio.EventPayload) {
		payld := kobrautils.NewMqttPayload("status", "query", nil)
		payld2 := kobrautils.NewMqttPayload("print", "update", M{"taskid": "0"})
		if err := mqtt.SendCommand(payld); err != nil {
			slog.Error(err.Error())
		}
		if err := mqtt.SendCommand(payld2); err != nil {
			slog.Error(err.Error())
		}
	})
}
