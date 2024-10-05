package server

import (
	"sync"

	"stefano.sonzogni/tic-tac-toe/internal/game"
)

type GameManager interface {
	HandleMessage(player, row, col int)
}

type Broadcaster interface {
	BroadcastGameState(gs game.GameState)
}

type ConcurrentGameManager struct {
	g           game.StatefulInteractableGame
	lock        sync.Mutex
	broadcaster Broadcaster
}

func NewConcurrentGameManager(g game.StatefulInteractableGame, broadcaster Broadcaster) ConcurrentGameManager {
	return ConcurrentGameManager{g, sync.Mutex{}, broadcaster}
}

func (gm *ConcurrentGameManager) HandleMessage(player, row, col int) {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	gm.g.PlaceMark(player, row, col)
	gm.broadcaster.BroadcastGameState(gm.g.State())
}
