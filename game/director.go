package game

import (
	"context"
	"fmt"
	"log"

	"github.com/JinWuZhao/sc2client"
	"github.com/JinWuZhao/sc2client/sc2proto"
	"github.com/golang/protobuf/proto"

	"github.com/jinwuzhao/stararena/command"
	"github.com/jinwuzhao/stararena/conf"
	"github.com/jinwuzhao/stararena/data"
	"github.com/jinwuzhao/stararena/msq"
	"github.com/jinwuzhao/stararena/state"
)

type Services struct {
	CmdQueue  *msq.Queue[command.Command]
	MsgQueue  *msq.Queue[string]
	GameState *state.Game
}

type Director struct {
	playerId     uint32
	sc2rpc       *sc2client.RpcClient
	cmdQueue     *msq.Queue[command.Command]
	msgQueue     *msq.Queue[string]
	gameState    *state.Game
	rankingSteps uint32
	gameEndStep  uint32
	gameEnding   bool
}

func NewDirector(cfg *conf.Conf, svc *Services) *Director {
	return &Director{
		cmdQueue:     svc.CmdQueue,
		msgQueue:     svc.MsgQueue,
		gameState:    svc.GameState,
		rankingSteps: cfg.RankingSteps,
	}
}

func (m *Director) OnStart(playerId uint32, rpc *sc2client.RpcClient) {
	log.Println("Director.OnStart():", playerId)
	m.playerId = playerId
	m.sc2rpc = rpc
	m.gameEndStep = 0
	m.gameEnding = false
	m.gameState.Prepare()
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
				log.Println("receive command:", chat.GetMessage())
				report, err := command.ParseReport(command.Context{}, chat.GetMessage())
				if err != nil {
					log.Println("command.ParseReport()", chat.GetMessage(), "error:", err)
				} else {
					if err := m.handleReport(report); err != nil {
						log.Println("m.handleReport()", report.String(), "error:", err)
					}
				}
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
	_, err := m.sc2rpc.Action(ctx, &sc2proto.RequestAction{
		Actions: sc2Actions,
	})
	if err != nil {
		log.Println("Director.OnStep(): m.sc2rpc.Action() error:", err)
		return
	}

	if !m.gameEnding && m.gameState.GetProgress() == state.GameProgressRanking {
		if m.gameEndStep == 0 {
			m.gameEndStep = st.Steps + m.rankingSteps
		} else if st.Steps >= m.gameEndStep {
			m.cmdQueue.Push((*command.EndGameCmd)(nil).New())
			m.gameEnding = true
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

func (m *Director) handleReport(report command.Command) error {
	switch report := report.(type) {
	case *command.DamageUnitReport:
		player := m.gameState.GetPlayer(report.Player)
		if player == nil {
			return fmt.Errorf("m.gameState.GetPlayer(): player %s not found", report.Player)
		}
		player.AddScore(int64(report.Damage))
	case *command.KillUnitReport:
		player := m.gameState.GetPlayer(report.Player)
		if player != nil {
			unit, ok := data.GetUnitByName(report.Unit)
			if !ok {
				return fmt.Errorf("data.GetUnitByName(): unit %s not found", report.Unit)
			}
			player.AddPoints(unit.Reward)
		}
	case *command.VictoryReport:
		m.gameState.Rank()
		for index, player := range m.gameState.GetRankedPlayers() {
			m.msgQueue.Push(fmt.Sprintf("第%d名：%s，得分：%d", index+1, player.GetName(), player.GetScore()))
		}
	}
	return nil
}
