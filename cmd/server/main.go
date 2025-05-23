package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"

	"github.com/OzkrOssa/ros-iface-streamer/internal/presentation"
	"github.com/OzkrOssa/ros-iface-streamer/pkg/config"
	"github.com/OzkrOssa/ros-iface-streamer/pkg/config/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	logger.Set(config.App)

	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			allowedOrigins := strings.Split(config.Http.AllowedOrigins, ",")
			return slices.Contains(allowedOrigins, r.Header.Get("Origin"))
		},
	}

	r := mux.NewRouter()
	r.HandleFunc("/ws", presentation.TrafficStreamerHandler(ctx, upgrader, config.Ros))

	listenAddr := fmt.Sprintf("%s:%s", config.Http.URL, config.Http.Port)

	go func() {
		slog.Info("Starting the HTTP server", "listen_address", listenAddr)
		if err := http.ListenAndServe(listenAddr, r); err != nil {
			slog.Error("Error starting the HTTP server", "error", err)
			return
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down the application")
}
