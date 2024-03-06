package websocket

import (
	"fmt"
	"log/slog"
	"time"

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

	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {
		// If ping received, send pong
		if string(ep.Data) == "ping" {
			ep.Kws.Emit([]byte("pong"), socketio.TextMessage)
			ep.Kws.Conn.SetReadDeadline(time.Now().Add(pongWait))
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
