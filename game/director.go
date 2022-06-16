package game

import (
	"context"
	"log"

	"github.com/JinWuZhao/sc2client"
	"github.com/JinWuZhao/sc2client/sc2proto"
	"google.golang.org/protobuf/proto"

	"github.com/jinwuzhao/stararena/command"
	"github.com/jinwuzhao/stararena/conf"
	"github.com/jinwuzhao/stararena/msq"
	"github.com/jinwuzhao/stararena/state"
)

type Services struct {
	CmdQueue  *msq.Queue[command.Command]
	MsgQueue  *msq.Queue[string]
	GameState *state.Game
}

type Director struct {
	config    *conf.Conf
	playerId  uint32
	sc2rpc    *sc2client.RpcClient
	cmdQueue  *msq.Queue[command.Command]
	msgQueue  *msq.Queue[string]
	gameState *state.Game
}

func NewDirector(cfg *conf.Conf, svc *Services) *Director {
	return &Director{
		config:    cfg,
		cmdQueue:  svc.CmdQueue,
		msgQueue:  svc.MsgQueue,
		gameState: svc.GameState,
	}
}

func (m *Director) OnStart(playerId uint32, rpc *sc2client.RpcClient) {
	log.Println("Director.OnStart():", playerId)
	m.playerId = playerId
	m.sc2rpc = rpc
	m.gameState.Start()
}

func (m *Director) OnStep(ctx context.Context, st *sc2client.StepState) {

chatLoop:
	for {
		select {
		case <-ctx.Done():
			break chatLoop
		case chat := <-st.ReceivedChats:
			if chat.GetPlayerId() == m.playerId {
				log.Println("chat:", chat.GetMessage())
			}
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

	if st.Steps%50 == 0 {
		m.gameState.HandleJoinPlayers()
		for _, player := range m.gameState.GetNewPlayers() {
			sc2Actions = append(sc2Actions, m.makeJoinAction(player))
		}
		m.gameState.ClearNewPlayers()
	}

	if len(sc2Actions) > 0 {
		_, err := m.sc2rpc.Action(ctx, &sc2proto.RequestAction{
			Actions: sc2Actions,
		})
		if err != nil {
			log.Println("Director.OnStep(): m.sc2rpc.Action() error:", err)
			return
		}
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
		}
		return nil
	case *command.SetUnitCmd:
		player := m.gameState.GetPlayer(cmd.Player())
		if player != nil {
			player.SetUnit(cmd.Unit())
		} else {
			return nil
		}
	default:
	}
	return &sc2proto.Action{
		ActionChat: &sc2proto.ActionChat{
			Channel: sc2proto.ActionChat_Team.Enum(),
			Message: proto.String(cmd.String()),
		},
	}
}

func (m *Director) makeJoinAction(player *state.Player) *sc2proto.Action {
	opts := []command.JoinGameOpts{
		command.JoinGameOptsPlayer(player.GetName(), player.GetSC2PlayerId()),
	}
	if player.IsBot() {
		opts = append(opts, command.JoinGameOptsBot())
	}
	cmd := (*command.JoinGameCmd)(nil).NewWithOpts(opts...)
	return &sc2proto.Action{
		ActionChat: &sc2proto.ActionChat{
			Channel: sc2proto.ActionChat_Team.Enum(),
			Message: proto.String(cmd.String()),
		},
	}
}
