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

func TrafficStreamerHandler(ctx context.Context, upgrader websocket.Upgrader, config *config.RouterOS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws := ws.New(upgrader)
		ws.Upgrader(w, r)

		mkt := mikrotik.New(config)
		manager := app.NewManager(mkt, ws)
		manager.HandleConnection(ctx)

	}
}
