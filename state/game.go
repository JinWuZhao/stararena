package state

import (
	"sort"
	"sync"

	"go.uber.org/atomic"
)

type GameProgress int32

const (
	GameProgressIdle = iota
	GameProgressPreparing
	GameProgressStarted
	GameProgressRanking
)

type Game struct {
	progress  atomic.Int32
	playerCap int
	players   map[string]*Player
	joinQueue chan *Player
	joinSet   map[string]struct{}
	m         sync.RWMutex
}

func NewGame(playerCap int, joinCap int) *Game {
	return &Game{
		playerCap: playerCap,
		players:   make(map[string]*Player, playerCap),
		joinQueue: make(chan *Player, joinCap),
		joinSet:   make(map[string]struct{}),
	}
}

func (m *Game) setProgress(progress GameProgress) {
	m.progress.Store(int32(progress))
}

func (m *Game) GetProgress() GameProgress {
	return GameProgress(m.progress.Load())
}

func (m *Game) Join(player *Player) bool {
	m.m.Lock()
	if _, ok := m.joinSet[player.GetName()]; ok {
		m.m.Unlock()
		return false
	}
	m.joinSet[player.GetName()] = struct{}{}
	m.m.Unlock()

	select {
	case m.joinQueue <- player:
		return true
	default:
		return false
	}
}

func (m *Game) Prepare() {
	m.setProgress(GameProgressPreparing)

	m.m.Lock()
	m.players = make(map[string]*Player, m.playerCap)
	m.m.Unlock()

joinLoop:
	for {
		select {
		case player := <-m.joinQueue:
			m.m.Lock()
			m.players[player.GetName()] = player
			delete(m.joinSet, player.GetName())
			m.m.Unlock()
		default:
			break joinLoop
		}
	}
}

func (m *Game) Start() {
	m.setProgress(GameProgressStarted)
}

func (m *Game) Rank() {
	m.setProgress(GameProgressRanking)
}

func (m *Game) Reset() {
	m.setProgress(GameProgressIdle)
}

func (m *Game) GetRankedPlayers() []*Player {
	m.m.RLock()
	defer m.m.RUnlock()

	players := make([]*Player, 0, m.playerCap)
	for _, p := range m.players {
		players = append(players, p)
	}
	sort.SliceStable(players, func(i, j int) bool {
		return players[i].GetScore() > players[j].GetScore()
	})
	return players
}

func (m *Game) GetPlayer(name string) *Player {
	m.m.RLock()
	defer m.m.RUnlock()

	return m.players[name]
}
