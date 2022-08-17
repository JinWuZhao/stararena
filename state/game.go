package state

import (
	"math/rand"
	"sync"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/atomic"

	"github.com/jinwuzhao/stararena/conf"
)

type GameProgress int32

const (
	GameProgressPreparing = iota
	GameProgressStarted
)

type Game struct {
	config         *conf.Conf
	progress       atomic.Int32
	players        map[string]*Player
	redPlayersNum  int
	bluePlayersNum int
	newPlayers     []string
	removePlayers  []string
	joinQueue      chan *Player
	joinSet        map[string]struct{}
	m              sync.RWMutex
	botIdGen       *snowflake.Node
}

func NewGame(config *conf.Conf) *Game {
	node, _ := snowflake.NewNode(1)
	return &Game{
		config:    config,
		players:   make(map[string]*Player, 0),
		joinQueue: make(chan *Player, config.JoinCap),
		joinSet:   make(map[string]struct{}),
		botIdGen:  node,
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

	m.m.Lock()
	for name := range m.players {
		delete(m.joinSet, name)
	}
	m.players = make(map[string]*Player, 0)
	m.redPlayersNum = 0
	m.bluePlayersNum = 0
	m.newPlayers = nil
	m.removePlayers = nil
	m.m.Unlock()

	m.setProgress(GameProgressStarted)
}

func (m *Game) HandleJoinPlayers() {
	for i := 0; i < 2; i++ {
		if m.GetProgress() != GameProgressStarted {
			break
		}

		if m.GetPlayerCount() >= m.config.PlayerCap &&
			m.FindOneBotPlayer(0) == nil {
			break
		}

		var player *Player
		select {
		case player = <-m.joinQueue:
			botPlayer := m.FindOneBotPlayer(player.GetSC2PlayerId())
			if botPlayer != nil {
				m.RemovePlayer(botPlayer.name)
			}
		default:
			m.m.RLock()
			if m.redPlayersNum != m.bluePlayersNum || m.redPlayersNum+m.bluePlayersNum < 10 {
				player = NewBotPlayer(m.botIdGen.Generate().String(), m.findAvailablePlayerId(false))
			}
			m.m.RUnlock()
		}
		if player == nil {
			break
		}

		m.m.Lock()
		if len(m.players) < m.config.PlayerCap {
			if m.redPlayersNum >= m.config.PlayerCap/2 && player.GetSC2PlayerId() == m.config.RedPlayer {
				player.SetSC2PlayerId(m.config.BluePlayer)
			} else if m.bluePlayersNum >= m.config.PlayerCap/2 && player.GetSC2PlayerId() == m.config.BluePlayer {
				player.SetSC2PlayerId(m.config.RedPlayer)
			}
			m.players[player.GetName()] = player
			m.newPlayers = append(m.newPlayers, player.GetName())
			if player.GetSC2PlayerId() == m.config.RedPlayer {
				m.redPlayersNum++
			} else if player.GetSC2PlayerId() == m.config.BluePlayer {
				m.bluePlayersNum++
			}
		}
		m.m.Unlock()
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

func (m *Game) GetPlayerCount() int {
	m.m.RLock()
	defer m.m.RUnlock()

	return len(m.players)
}

func (m *Game) RemovePlayer(name string) *Player {
	m.m.Lock()
	defer m.m.Unlock()

	player := m.players[name]
	delete(m.players, name)
	delete(m.joinSet, name)
	if player.GetSC2PlayerId() == m.config.RedPlayer {
		m.redPlayersNum--
	} else if player.GetSC2PlayerId() == m.config.BluePlayer {
		m.bluePlayersNum--
	}
	m.removePlayers = append(m.removePlayers, name)
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

func (m *Game) GetRemovePlayers() []string {
	m.m.RLock()
	defer m.m.RUnlock()

	var players []string
	for _, name := range m.removePlayers {
		players = append(players, name)
	}
	return players
}

func (m *Game) ClearNewPlayers() {
	m.m.Lock()
	defer m.m.Unlock()

	m.newPlayers = nil
}

func (m *Game) ClearRemovePlayers() {
	m.m.Lock()
	defer m.m.Unlock()

	m.removePlayers = nil
}

func (m *Game) FindOneBotPlayer(sc2PlayerId uint32) *Player {
	m.m.RLock()
	defer m.m.RUnlock()

	var bot *Player
	for _, p := range m.players {
		if p.IsBot() {
			bot = p
			if sc2PlayerId == 0 || p.GetSC2PlayerId() == sc2PlayerId {
				break
			}
		}
	}
	return bot
}

func (m *Game) findAvailablePlayerId(humanOnly bool) uint32 {
	var redPlayersNum, bluePlayersNum int
	if humanOnly {
		for _, p := range m.players {
			if !p.IsBot() {
				if p.GetSC2PlayerId() == m.config.RedPlayer {
					redPlayersNum += 1
				} else if p.GetSC2PlayerId() == m.config.BluePlayer {
					bluePlayersNum += 1
				}
			}
		}
	} else {
		redPlayersNum = m.redPlayersNum
		bluePlayersNum = m.bluePlayersNum
	}

	if redPlayersNum > bluePlayersNum {
		return m.config.BluePlayer
	}
	if redPlayersNum < bluePlayersNum {
		return m.config.RedPlayer
	}
	return []uint32{m.config.BluePlayer, m.config.RedPlayer}[rand.Intn(2)]
}

func (m *Game) FindAvailablePlayerId(humanOnly bool) uint32 {
	m.m.RLock()
	defer m.m.RUnlock()
	return m.findAvailablePlayerId(humanOnly)
}
