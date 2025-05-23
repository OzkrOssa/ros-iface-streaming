package presentation

import (
	"context"
	"net/http"

	"github.com/OzkrOssa/ros-iface-streamer/internal/app"
	"github.com/OzkrOssa/ros-iface-streamer/internal/infra/mikrotik"
	"github.com/OzkrOssa/ros-iface-streamer/internal/infra/ws"
	"github.com/OzkrOssa/ros-iface-streamer/pkg/config"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TrafficStreamerHandler(ctx context.Context, config *config.RouterOS) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws := ws.New(upgrader)
		ws.Upgrader(w, r)

		mkt := mikrotik.New(config)
		manager := app.NewManager(mkt, ws)
		manager.HandleConnection(ctx)
	}
}
