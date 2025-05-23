package domain

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WsClient interface {
	Upgrader(w http.ResponseWriter, r *http.Request) error
	AddClient(value any) map[*websocket.Conn]any
	DeleteClient(key *websocket.Conn)
	CurrentConnection() *websocket.Conn
	ReadMessage() ([]byte, error)
	WriteMessage([]byte) error
	WriteErrorMessage(message []byte) error
	WriteJSON(any) error
	Close() error
}

type WsPayload struct {
	Host  string `json:"host"`
	Iface string `json:"iface"`
}
