package game

import (
	"context"
	"log"

	"github.com/JinWuZhao/sc2client"
	"github.com/JinWuZhao/sc2client/sc2proto"
	"github.com/golang/protobuf/proto"

	"github.com/jinwuzhao/stararena/command"
	"github.com/jinwuzhao/stararena/msq"
	"github.com/jinwuzhao/stararena/state"
)

type Services struct {
	CmdQueue  *msq.Queue[command.Command]
	GameState *state.Game
}

type Director struct {
	playerId  uint32
	sc2rpc    *sc2client.RpcClient
	cmdQueue  *msq.Queue[command.Command]
	gameState *state.Game
}

func NewDirector(svc *Services) *Director {
	return &Director{
		cmdQueue:  svc.CmdQueue,
		gameState: svc.GameState,
	}
}

func (m *Director) OnStart(playerId uint32, rpc *sc2client.RpcClient) {
	log.Println("Director.OnStart():", playerId)
	m.playerId = playerId
	m.sc2rpc = rpc
	m.gameState.Prepare()
	m.gameState.Start()
}

func (m *Director) OnStep(ctx context.Context, state *sc2client.StepState) {
chatLoop:
	for {
		select {
		case <-ctx.Done():
			break chatLoop
		case chat := <-state.ReceivedChats:
			if chat.GetPlayerId() == m.playerId {
				log.Println("receive command:", chat.GetMessage())
			}
			// TODO 处理消息中的命令
		default:
			break chatLoop
		}
	}

	var sc2Actions []*sc2proto.Action
cmdLoop:
	for {
		select {
		case <-ctx.Done():
			break cmdLoop
		case cmd := <-m.cmdQueue.PopChan():
			action := m.handleCommand(cmd)
			if action != nil {
				sc2Actions = append(sc2Actions, action)
			}
		default:
			break cmdLoop
		}
	}
	_, err := m.sc2rpc.Action(ctx, &sc2proto.RequestAction{
		Actions: sc2Actions,
	})
	if err != nil {
		log.Println("Director.OnStep(): m.sc2rpc.Action() error:", err)
		return
	}
}

func (m *Director) OnEnd(_ sc2proto.Result) {
	m.gameState.Reset()
}

func (m *Director) handleCommand(cmd command.Command) *sc2proto.Action {
	switch cmd := cmd.(type) {
	case *command.JoinGameCmd:
		if m.gameState.Join(state.NewPlayer(cmd.Player(), cmd.SC2PlayerId())) {
			log.Println(cmd.Player(), "joined game player", cmd.SC2PlayerId())
		} else {
			// TODO show message
		}
		return nil
	default:
		return &sc2proto.Action{
			ActionChat: &sc2proto.ActionChat{
				Channel: sc2proto.ActionChat_Team.Enum(),
				Message: proto.String(cmd.String()),
			},
		}
	}
}
