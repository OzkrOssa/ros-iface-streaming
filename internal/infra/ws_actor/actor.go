package wsactor

import (
	"log"

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
	switch msg := ctx.Message().(type) {
	case *domain.Traffic:
		err := w.wsclient.WriteJSON(
			map[string]any{
				"iface": msg.Iface,
				"rx":    msg.Rx,
				"tx":    msg.Tx,
			})

		if err != nil {
			log.Println(err)
		}
	}
}
