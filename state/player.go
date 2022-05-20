package state

import (
	"go.uber.org/atomic"
)

type Player struct {
	name        atomic.String // 玩家名字
	sc2PlayerId atomic.Uint32 // SC2玩家ID
	unit        atomic.String // 操作单位
}

func NewPlayer(name string, sc2PlayerId uint32) *Player {
	player := new(Player)
	player.name.Store(name)
	player.sc2PlayerId.Store(sc2PlayerId)
	return player
}

func (m *Player) SetName(name string) {
	m.name.Store(name)
}

func (m *Player) GetName() string {
	return m.name.Load()
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
