package server

import (
	"sync"

	"stefano.sonzogni/tic-tac-toe/internal/game"
)

type GameManager interface {
	HandleMessage(player, row, col int)
}

type ConcurrentGameManager struct {
	g    game.StatefulInteractableGame
	lock sync.Mutex
}

func (gm *ConcurrentGameManager) HandleMessage(player, row, col int) {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	gm.g.PlaceMark(player, row, col)
}
