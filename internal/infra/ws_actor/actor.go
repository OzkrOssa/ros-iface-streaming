package wsactor

import (
	"log/slog"

	"github.com/OzkrOssa/ros-iface-streamer/internal/domain"
	"github.com/anthdm/hollywood/actor"
)

type wsActor struct {
	wsclient domain.WsClient
}

func New(wsclient domain.WsClient) actor.Producer {
	return func() actor.Receiver {
		return &wsActor{wsclient: wsclient}
	}
}

func (w *wsActor) Receive(ctx *actor.Context) {
	remote := w.wsclient.CurrentConnection().RemoteAddr().String()
	actorID := ctx.PID().ID
	switch msg := ctx.Message().(type) {
	case actor.Started:
		slog.Info("websocket actor started listener", "remote", remote, "actor", actorID)
	case actor.Stopped:
		slog.Info("websocket actor stopped listener", "remote", remote, "actor", actorID)
	case *domain.Traffic:
		err := w.wsclient.WriteJSON(
			map[string]any{
				"iface": msg.Iface,
				"rx":    msg.Rx,
				"tx":    msg.Tx,
			})

		if err != nil {
			slog.Error("Error writing JSON", "error", err)
			return
		}
	}
}
