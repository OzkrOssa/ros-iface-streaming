package ws

import (
	"net/http"
	"sync"

	"github.com/OzkrOssa/ros-iface-streamer/internal/domain"
	"github.com/gorilla/websocket"
)

type WsClient struct {
	upgrader websocket.Upgrader
	conn     *websocket.Conn
	client   map[*websocket.Conn]any
	mu       sync.Mutex
}

var _ domain.WsClient = (*WsClient)(nil)

func New(u websocket.Upgrader) *WsClient {
	return &WsClient{
		upgrader: u,
		client:   make(map[*websocket.Conn]any),
	}
}

func (c *WsClient) Upgrader(w http.ResponseWriter, r *http.Request) error {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *WsClient) AddClient(value any) map[*websocket.Conn]any {
	c.mu.Lock()
	c.client[c.conn] = value
	c.mu.Unlock()
	return c.client
}

func (c *WsClient) DeleteClient(key *websocket.Conn) {
	c.mu.Lock()
	delete(c.client, key)
	c.mu.Unlock()
}

func (c *WsClient) CurrentConnection() *websocket.Conn {
	return c.conn
}

func (c *WsClient) ReadMessage() ([]byte, error) {
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *WsClient) WriteMessage(message []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, message)
}

func (c *WsClient) WriteErrorMessage(message []byte) error {
	return c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseUnsupportedData, string(message)))
}

func (c *WsClient) WriteJSON(message any) error {
	return c.conn.WriteJSON(message)
}

func (c *WsClient) Close() error {
	return c.conn.Close()
}
