package app

import (
	"context"
	"encoding/json"

	"github.com/OzkrOssa/ros-iface-streamer/internal/domain"
	wsactor "github.com/OzkrOssa/ros-iface-streamer/internal/infra/ws_actor"
	"github.com/anthdm/hollywood/actor"
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
		return err
	}

	var payload domain.WsPayload
	if err = json.Unmarshal(message, &payload); err != nil {
		_ = m.ws.WriteErrorMessage([]byte("invalid playload"))
		_ = m.ws.Close()
		return err
	}

	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		return err
	}

	wsactors := wsactor.New(m.ws)
	pid := engine.Spawn(wsactors, "ws-actor")

	m.ws.AddClient(pid)

	if err := m.mkt.Connect(payload.Host); err != nil {
		_ = m.ws.Close()
		engine.Stop(pid)
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
			cancel()
			return
		}

		for t := range traffic {
			engine.Send(pid, t)
		}
	}()

	<-ctx.Done()
	engine.Stop(pid)
	m.ws.DeleteClient(m.ws.CurrentConnection())
	_ = m.ws.Close()
	return nil
}
