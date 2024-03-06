package websocket

import (
	"fmt"
	"log/slog"

	"github.com/anjomro/kobra-unleashed/server/kobrautils"
	"github.com/anjomro/kobra-unleashed/server/mqtt"
	"github.com/gofiber/contrib/socketio"
)

type M map[string]interface{}

// map of clients uuid
var clients = make(map[string]*socketio.Websocket)

func SetupWebsocket() {

	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		// Add client id to clinets
		clients[ep.Kws.UUID] = ep.Kws
		fmt.Printf("Client %s connected\n", ep.Kws.UUID)
	})

	socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
		// Remove client id from clients
		delete(clients, ep.Kws.UUID)
		fmt.Printf("Client %s closed\n", ep.Kws.UUID)
	})

	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		ep.Kws.Close()
		// Remove client id from clients
		delete(clients, ep.Kws.UUID)
		fmt.Printf("Client %s disconnected\n", ep.Kws.UUID)
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
