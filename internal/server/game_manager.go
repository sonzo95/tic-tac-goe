package server

import (
	"sync"

	"stefano.sonzogni/tic-tac-toe/internal/game"
)

type GameManager interface {
	HandleMessage(player, row, col int)
}

type ConcurrentGameManager struct {
	g           game.StatefulInteractableGame
	lock        sync.Mutex
	broadcaster Broadcaster
}

func NewConcurrentGameManager(broadcaster Broadcaster) ConcurrentGameManager {
	g := game.NewGame()
	return ConcurrentGameManager{&g, sync.Mutex{}, broadcaster}
}

func (gm *ConcurrentGameManager) Start() {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	gm.broadcaster.BroadcastGameState(gm.g.State())
}

func (gm *ConcurrentGameManager) HandleMessage(player, row, col int) {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	err := gm.g.PlaceMark(player, row, col)
	if err == nil {
		gm.broadcaster.BroadcastGameState(gm.g.State())
	}
}
