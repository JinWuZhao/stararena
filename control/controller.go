package control

import (
	"context"
	"fmt"
	"log"

	"github.com/JinWuZhao/bilidanmu"

	"github.com/jinwuzhao/stararena/command"
	"github.com/jinwuzhao/stararena/conf"
	"github.com/jinwuzhao/stararena/msq"
	"github.com/jinwuzhao/stararena/state"
)

type Services struct {
	CmdQueue  *msq.Queue[command.Command]
	GameState *state.Game
}

type Controller struct {
	config    *conf.Conf
	client    *bilidanmu.Client
	cmdQueue  *msq.Queue[command.Command]
	gameState *state.Game
}

func NewController(cfg *conf.Conf, svc *Services) (*Controller, error) {
	client, err := bilidanmu.NewClient(cfg.RoomId)
	if err != nil {
		return nil, fmt.Errorf("bilidanmu.NewClient() error: %w", err)
	}
	return &Controller{
		config:    cfg,
		client:    client,
		cmdQueue:  svc.CmdQueue,
		gameState: svc.GameState,
	}, nil
}

func (s *Controller) Start(ctx context.Context) error {
	err := s.client.Start(ctx, s.ReceiveMessage)
	if err != nil {
		return fmt.Errorf("s.client.Start() error: %w", err)
	}
	return nil
}

func (s *Controller) ReceiveMessage(message bilidanmu.Message) {
	switch m := message.(type) {
	case *bilidanmu.DanMuMsg:
		log.Printf("%d-%s | %d-%s: %s\n", m.MedalLevel, m.MedalName, m.ULevel, m.Uname, m.Text)
		ctx := command.Context{
			SC2RedPlayer:  s.config.RedPlayer,
			SC2BluePlayer: s.config.BluePlayer,
			Player:        m.Uname,
		}
		player := s.gameState.GetPlayer(m.Uname)
		if player != nil {
			ctx.Unit = player.GetUnit()
		}
		cmd, err := command.ParseCommand(ctx, m.Text)
		if err == nil {
			if !s.cmdQueue.Push(cmd) {
				log.Println("[WARN] command queue overflowed")
			}
		}
	case *bilidanmu.Gift:
		log.Printf("%s %s 价值 %d 的 %s\n", m.UUname, m.Action, m.Price, m.GiftName)
	}
}

func (s *Controller) WaitForStop() {
	err := s.client.WaitForStop()
	if err != nil {
		log.Println("s.client.WaitForStop() error:", err)
	}
}
