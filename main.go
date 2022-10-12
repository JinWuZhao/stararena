package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"time"

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

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	cfg, err := conf.NewConf(*confPath)
	if err != nil {
		log.Println("conf.NewConf() error:", err)
		return
	}
	if cfg.Log.FilePath != "" {
		f, err := os.OpenFile(cfg.Log.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Println("os.OpenFile(", cfg.Log.FilePath, ") error:", err)
			return
		}
		defer f.Close()
		log.SetOutput(f)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cmdQueue := msq.NewQueue[command.Command](cfg.CmdQueueCap)
	msgQueue := msq.NewQueue[string](cfg.MsgQueueCap)
	gameState := state.NewGame(cfg)

	var controller *control.Controller
	if cfg.BiliDanMu.Enable {
		controller, err = control.NewController(cfg,
			&control.Services{
				CmdQueue:  cmdQueue,
				GameState: gameState,
			})
		if err != nil {
			log.Println("control.NewController() error:", err)
			return
		}
	} else {
		controller = control.NewFakeController(cfg,
			&control.Services{
				CmdQueue:  cmdQueue,
				GameState: gameState,
			})
	}

	err = controller.Start(ctx)
	if err != nil {
		log.Println("controller.Start() error:", err)
		return
	}
	defer controller.WaitForStop()

	director := game.NewDirector(
		cfg,
		&game.Services{
			CmdQueue:  cmdQueue,
			GameState: gameState,
			MsgQueue:  msgQueue,
		})
	audience := game.NewAudience(msgQueue)

	gameMaps := make([]sc2client.GameMap, 0, len(cfg.GameMaps))
	for _, m := range cfg.GameMaps {
		gameMaps = append(gameMaps, sc2client.GameMap{
			Name:       filepath.Base(m),
			SourcePath: m,
		})
	}

	err = sc2client.RunGame(ctx,
		gameMaps,
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
}
