package state

import "go.uber.org/atomic"

type Player struct {
	name        string        // 玩家名字
	isBot       bool          // 是否为机器人
	sc2PlayerId atomic.Uint32 // SC2玩家ID
	unit        atomic.String // 操作单位
}

func NewPlayer(name string, sc2PlayerId uint32) *Player {
	player := new(Player)
	player.name = name
	player.sc2PlayerId.Store(sc2PlayerId)
	return player
}

func NewBotPlayer(name string, sc2PlayerId uint32) *Player {
	player := NewPlayer(name, sc2PlayerId)
	player.isBot = true
	return player
}

func (m *Player) GetName() string {
	return m.name
}

func (m *Player) IsBot() bool {
	return m.isBot
}

func (m *Player) SetSC2PlayerId(playerId uint32) {
	m.sc2PlayerId.Store(playerId)
}

func (m *Player) GetSC2PlayerId() uint32 {
	return m.sc2PlayerId.Load()
}

func (m *Player) SetUnit(unit string) {
	m.unit.Store(unit)
}

func (m *Player) GetUnit() string {
	return m.unit.Load()
}
