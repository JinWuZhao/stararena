package control

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

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
	config         *conf.Conf
	client         *bilidanmu.Client
	cmdQueue       *msq.Queue[command.Command]
	gameState      *state.Game
	fakeClientStop chan struct{}
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

func NewFakeController(cfg *conf.Conf, svc *Services) *Controller {
	return &Controller{
		config:         cfg,
		cmdQueue:       svc.CmdQueue,
		gameState:      svc.GameState,
		fakeClientStop: make(chan struct{}, 1),
	}
}

var fakeMsgRegex = regexp.MustCompile(`^([a-zA-Z0-9\p{Han}_.-]{1,10}):\s*([a-zA-Z0-9]+)$`)

func (s *Controller) Start(ctx context.Context) error {
	if s.client != nil {
		err := s.client.Start(ctx, s.ReceiveMessage)
		if err != nil {
			return fmt.Errorf("s.client.Start() error: %w", err)
		}
	} else {
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			for ctx.Err() == nil && scanner.Scan() {
				matches := fakeMsgRegex.FindStringSubmatch(scanner.Text())
				if len(matches) >= 3 {
					msg := new(bilidanmu.DanMuMsg)
					msg.Uname = matches[1]
					msg.Text = matches[2]
					s.ReceiveMessage(msg)
				}
			}
			if err := scanner.Err(); err != nil {
				log.Println("scanner.Scan() error:", err)
			}
			s.fakeClientStop <- struct{}{}
		}()
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
			PlayerUID:     m.UID,
			Streamer:      s.config.Streamer,
			State:         s.gameState,
		}
		cmd, err := command.ParseCommand(ctx, m.Text)
		if err == nil {
			if !s.cmdQueue.Push(cmd) {
				log.Println("[WARN] command queue overflowed")
			}
		}
	case *bilidanmu.Gift:
		log.Printf("%s %s 价值 %d 的 %s X%d\n", m.Uname, m.Action, m.Price*m.Number, m.GiftName, m.Number)
		cmd := (*command.GiftItemCmd)(nil).NewWithOpts(command.GiftItemOptsGift(m.Uname, m.GiftName, m.Number))
		if !s.cmdQueue.Push(cmd) {
			log.Println("[ERROR] command queue overflowed for gift")
		}
	}
}

func (s *Controller) WaitForStop() {
	if s.client != nil {
		err := s.client.WaitForStop()
		if err != nil {
			log.Println("s.client.WaitForStop() error:", err)
		}
	} else {
		<-s.fakeClientStop
		log.Println("fake controller stopped")
	}
}
