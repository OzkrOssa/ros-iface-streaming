package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/OzkrOssa/ros-iface-streamer/internal/presentation"
	"github.com/OzkrOssa/ros-iface-streamer/pkg/config"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config, err := config.New()
	if err != nil {
		log.Println(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/ws", presentation.TrafficStreamerHandler(ctx, config.Ros))

	go func() {
		log.Println("🌐 Servidor HTTP iniciado en :8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatalf("❌ Error en servidor HTTP: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("🛑 Finalizando conexiones...")
}
