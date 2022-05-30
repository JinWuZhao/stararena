package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"

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
var mapsPath = flag.String("m", "sc2maps/product", "map files path")

func main() {
	cfg, err := conf.NewConf(*confPath)
	if err != nil {
		log.Println("conf.NewConf() error:", err)
		return
	}

	if *mapsPath != "" {
		files, err := os.ReadDir(*mapsPath)
		if err == nil {
			sc2Path, err := sc2client.GetSC2InstallDir()
			if err != nil {
				log.Println("sc2client.GetSC2InstallDir() error:", err)
				return
			}
			sc2MapPath := filepath.Join(sc2Path, "Maps")
			if _, err := os.Stat(sc2MapPath); os.IsNotExist(err) {
				err = os.Mkdir(sc2MapPath, os.ModePerm)
				if err != nil {
					log.Println("os.Mkdir() error:", sc2MapPath, err)
					return
				}
			}
			for _, f := range files {
				if !f.IsDir() {
					srcPath := filepath.Join(*mapsPath, f.Name())
					content, err := os.ReadFile(srcPath)
					if err != nil {
						log.Println("failed to read map file:", srcPath, err)
						return
					}
					dstPath := filepath.Join(sc2MapPath, f.Name())
					err = os.WriteFile(dstPath, content, os.ModePerm)
					if err != nil {
						log.Println("failed to write map file:", dstPath, err)
						return
					}
				}
			}
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cmdQueue := msq.NewQueue[command.Command](cfg.CmdQueueCap)
	msgQueue := msq.NewQueue[string](cfg.MsgQueueCap)
	gameState := state.NewGame(cfg)

	if cfg.BiliDanMu.Enable {
		controller, err := control.NewController(cfg, &control.Services{
			CmdQueue:  cmdQueue,
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
		defer controller.WaitForStop()
	}

	director := game.NewDirector(
		cfg,
		&game.Services{
			CmdQueue:  cmdQueue,
			GameState: gameState,
		})
	audience := game.NewAudience(msgQueue)

	err = sc2client.RunGame(ctx,
		cfg.GameMaps,
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
