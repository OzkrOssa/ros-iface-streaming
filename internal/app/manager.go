package app

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/OzkrOssa/ros-iface-streamer/internal/domain"
	"github.com/google/uuid"
)

type Manager struct {
	mkt domain.Mikrotik
	ws  domain.WsClient
}

func NewManager(
	mkt domain.Mikrotik,
	ws domain.WsClient,
) *Manager {
	return &Manager{
		mkt: mkt,
		ws:  ws,
	}
}

func (m *Manager) HandleConnection(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	message, err := m.ws.ReadMessage()
	if err != nil {
		_ = m.ws.Close()
		slog.Error("Error reading message", "error", err)
		return err
	}

	var payload domain.WsPayload
	if err = json.Unmarshal(message, &payload); err != nil {
		_ = m.ws.WriteErrorMessage([]byte("invalid playload"))
		_ = m.ws.Close()
		slog.Error("Error unmarshalling message", "error", err)
		return err
	}

	if payload.Host == "" || payload.Iface == "" {
		_ = m.ws.WriteErrorMessage([]byte("invalid playload"))
		_ = m.ws.Close()
		slog.Error("Invalid playload", "payload", payload)
		return err
	}

	clientConn := m.ws.AddClient(uuid.New().String())
	slog.Info("New connection", "remote", m.ws.CurrentConnection().RemoteAddr(), "connection", clientConn)

	if err := m.mkt.Connect(payload.Host); err != nil {
		_ = m.ws.Close()
		slog.Error("Error connecting to mikrotik", "error", err)
		return err
	}
	defer m.mkt.Close()

	go func() {
		defer cancel()
		for {
			_, err := m.ws.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	go func() {
		traffic, err := m.mkt.GetStreamingTraffic(ctx, payload.Iface)
		if err != nil {
			slog.Error("Error getting streaming traffic", "error", err)
			cancel()
			return
		}

		for t := range traffic {
			err := m.ws.WriteJSON(
				map[string]any{
					"iface": t.Iface,
					"rx":    t.Rx,
					"tx":    t.Tx,
				})
			if err != nil {
				slog.Error("Error writing JSON", "error", err)
				return
			}
		}
	}()

	<-ctx.Done()
	m.ws.DeleteClient(m.ws.CurrentConnection())
	_ = m.ws.Close()
	slog.Info("Connection closed", "remote", m.ws.CurrentConnection().RemoteAddr())
	return nil
}
