package state

import (
	"go.uber.org/atomic"
)

type Player struct {
	name        atomic.String // 玩家名字
	sc2PlayerId atomic.Uint32 // SC2玩家ID
	score       atomic.Int64  // 分数
	unit        atomic.String // 操作单位
	points      atomic.Int64  // 消费点数
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

func (m *Player) SetScore(score int64) {
	m.score.Store(score)
}

func (m *Player) AddScore(score int64) int64 {
	return m.score.Add(score)
}

func (m *Player) GetScore() int64 {
	return m.score.Load()
}

func (m *Player) SetUnit(unit string) {
	m.unit.Store(unit)
}

func (m *Player) GetUnit() string {
	return m.unit.Load()
}

func (m *Player) SetPoints(points int64) {
	m.points.Store(points)
}

func (m *Player) AddPoints(points int64) int64 {
	return m.points.Add(points)
}

func (m *Player) GetPoints() int64 {
	return m.points.Load()
}
