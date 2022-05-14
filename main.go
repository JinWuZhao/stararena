package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/JinWuZhao/sc2client"
	"github.com/JinWuZhao/sc2client/sc2proto"

	"github.com/jinwuzhao/stararena/command"
	"github.com/jinwuzhao/stararena/conf"
	"github.com/jinwuzhao/stararena/control"
	"github.com/jinwuzhao/stararena/game"
	"github.com/jinwuzhao/stararena/msq"
	"github.com/jinwuzhao/stararena/state"
)

var confPath = flag.String("c", "conf/conf.toml", "config file path")

func main() {
	cfg, err := conf.NewConf(*confPath)
	if err != nil {
		log.Println("conf.NewConf() error:", err)
		return
	}

	cmdQueue := msq.NewQueue[command.Command](cfg.CmdQueueCap)
	msgQueue := msq.NewQueue[string](cfg.MsgQueueCap)
	gameState := state.NewGame(cfg.PlayerCap, cfg.JoinCap)
	director := game.NewDirector(&game.Services{
		CmdQueue:  cmdQueue,
		GameState: gameState,
	})
	audience := game.NewAudience(msgQueue)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	controller, err := control.NewController(cfg, &control.Services{
		CmdQueue:  cmdQueue,
		MsgQueue:  msgQueue,
		GameState: gameState,
	})
	if err != nil {
		log.Println("control.NewController() error:", err)
		return
	}

	err = controller.Start(ctx)
	if err != nil {
		log.Println("controller.Start() error:", err)
		return
	}

	err = sc2client.RunGame(ctx,
		cfg.GameMap,
		[]*sc2client.PlayerSetup{
			{
				Type:       sc2proto.PlayerType_Participant,
				Race:       sc2proto.Race_Random,
				Name:       cfg.DirectorName,
				Difficulty: sc2proto.Difficulty_Easy,
				AIBuild:    sc2proto.AIBuild_RandomBuild,
				Agent:      director,
			},
			{
				Type:       sc2proto.PlayerType_Participant,
				Race:       sc2proto.Race_Random,
				Name:       cfg.AudienceName,
				Difficulty: sc2proto.Difficulty_Easy,
				AIBuild:    sc2proto.AIBuild_RandomBuild,
				Agent:      audience,
			}},
		true)
	if err != nil {
		log.Println("sc2client.RunGame() error:", err)
		return
	}

	controller.WaitForStop()
}
