package game

import (
	"context"

	"github.com/JinWuZhao/sc2client"
	"github.com/JinWuZhao/sc2client/sc2proto"

	"github.com/jinwuzhao/stararena/msq"
)

type Audience struct {
	rpc  *sc2client.RpcClient
	cmdQ *msq.CommandQueue
}

func NewAudience(cmdQ *msq.CommandQueue) *Audience {
	return &Audience{
		cmdQ: cmdQ,
	}
}

func (m *Audience) OnStart(playerId uint32, rpc *sc2client.RpcClient) {
	m.rpc = rpc
}

func (m *Audience) OnStep(ctx context.Context, state *sc2client.StepState) {
	//TODO implement me
	panic("implement me")
}

func (m *Audience) OnEnd(result sc2proto.Result) {
	//TODO implement me
	panic("implement me")
}
