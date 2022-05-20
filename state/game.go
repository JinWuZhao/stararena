package state

import (
	"sync"

	"go.uber.org/atomic"
)

type GameProgress int32

const (
	GameProgressPreparing = iota
	GameProgressStarted
)

type Game struct {
	progress   atomic.Int32
	playerCap  int
	players    map[string]*Player
	newPlayers []string
	joinQueue  chan *Player
	joinSet    map[string]struct{}
	m          sync.RWMutex
}

func NewGame(playerCap int, joinCap int) *Game {
	return &Game{
		playerCap: playerCap,
		players:   make(map[string]*Player, 0),
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
	m.m.RLock()
	if _, ok := m.joinSet[player.GetName()]; ok {
		m.m.RUnlock()
		return false
	}
	m.m.RUnlock()

	select {
	case m.joinQueue <- player:
		m.m.Lock()
		if _, ok := m.joinSet[player.GetName()]; !ok {
			m.joinSet[player.GetName()] = struct{}{}
		}
		m.m.Unlock()
		return true
	default:
		return false
	}
}

func (m *Game) Start() {
	if m.GetProgress() != GameProgressPreparing {
		return
	}

	m.setProgress(GameProgressStarted)

	m.m.Lock()
	for name := range m.players {
		delete(m.joinSet, name)
	}
	m.players = make(map[string]*Player, 0)
	m.newPlayers = nil
	m.m.Unlock()
}

func (m *Game) Step() {
	for i := 0; i < 10; i++ {
		if m.GetProgress() != GameProgressStarted {
			break
		}

		m.m.RLock()
		if len(m.players) >= m.playerCap {
			m.m.RUnlock()
			break
		}
		m.m.RUnlock()

		select {
		case player := <-m.joinQueue:
			m.m.Lock()
			if len(m.players) < m.playerCap {
				m.players[player.GetName()] = player
				m.newPlayers = append(m.newPlayers, player.GetName())
			}
			m.m.Unlock()
		default:
		}
	}
}

func (m *Game) Reset() {
	m.setProgress(GameProgressPreparing)
}

func (m *Game) GetPlayer(name string) *Player {
	m.m.RLock()
	defer m.m.RUnlock()

	return m.players[name]
}

func (m *Game) RemovePlayer(name string) *Player {
	m.m.Lock()
	defer m.m.Unlock()

	player := m.players[name]
	delete(m.players, name)
	delete(m.joinSet, name)
	return player
}

func (m *Game) GetNewPlayers() []*Player {
	m.m.RLock()
	defer m.m.RUnlock()

	var players []*Player
	for _, name := range m.newPlayers {
		if p, ok := m.players[name]; ok {
			players = append(players, p)
		}
	}
	return players
}

func (m *Game) ClearNewPlayers() {
	m.m.Lock()
	defer m.m.Unlock()

	m.newPlayers = nil
}
