package game

import (
	"context"
	"log"

	"github.com/JinWuZhao/sc2client"
	"github.com/JinWuZhao/sc2client/sc2proto"
	"google.golang.org/protobuf/proto"

	"github.com/jinwuzhao/stararena/msq"
)

type Audience struct {
	sc2rpc   *sc2client.RpcClient
	msgQueue *msq.Queue[string]
}

func NewAudience(msgQueue *msq.Queue[string]) *Audience {
	return &Audience{
		msgQueue: msgQueue,
	}
}

func (m *Audience) OnStart(playerId uint32, rpc *sc2client.RpcClient) {
	m.sc2rpc = rpc
}

func (m *Audience) OnStep(ctx context.Context, state *sc2client.StepState) {
	var sc2Actions []*sc2proto.Action
msgLoop:
	for {
		select {
		case <-ctx.Done():
			break msgLoop
		case msg := <-m.msgQueue.PopChan():
			action := &sc2proto.Action{
				ActionChat: &sc2proto.ActionChat{
					Channel: sc2proto.ActionChat_Broadcast.Enum(),
					Message: proto.String(msg),
				},
			}
			sc2Actions = append(sc2Actions, action)
		default:
			break msgLoop
		}
	}

	_, err := m.sc2rpc.Action(ctx, &sc2proto.RequestAction{
		Actions: sc2Actions,
	})
	if err != nil {
		log.Println("Audience.OnStep(): m.sc2rpc.Action() error:", err)
		return
	}
}

func (m *Audience) OnEnd(result sc2proto.Result) {
}
