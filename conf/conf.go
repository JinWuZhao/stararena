package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type BiliDanMu struct {
	RoomId uint32 `toml:"room_id"`
}

type SC2Client struct {
	RedPlayer    uint32 `toml:"red_player"`
	BluePlayer   uint32 `toml:"blue_player"`
	GameMap      string `toml:"game_map"`
	DirectorName string `toml:"director_name"`
	AudienceName string `toml:"audience_name"`
}

type Msq struct {
	CmdQueueCap int `toml:"cmd_queue_cap"`
	MsgQueueCap int `toml:"msg_queue_cap"`
}

type State struct {
	PlayerCap int `toml:"player_cap"`
	JoinCap   int `toml:"join_cap"`
}

type Conf struct {
	BiliDanMu `toml:"bilidanmu"`
	SC2Client `toml:"sc2client"`
	Msq       `toml:"msq"`
	State     `toml:"state"`
}

func NewConf(file string) (*Conf, error) {
	c := new(Conf)
	_, err := toml.DecodeFile(file, c)
	if err != nil {
		return nil, fmt.Errorf("toml.DecodeFile() error: %w", err)
	}
	return c, nil
}
