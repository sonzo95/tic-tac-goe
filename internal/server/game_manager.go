package server

import (
	"stefano.sonzogni/tic-tac-toe/internal/game"
)

type ConcurrentGameManager struct {
	g      game.StatefulInteractableGame
	p1, p2 *player
}

func NewConcurrentGameManager(p1, p2 *player) *ConcurrentGameManager {
	g := game.NewGame()
	return &ConcurrentGameManager{&g, p1, p2}
}

func (gm *ConcurrentGameManager) Start() {
	gs := gm.g.State()
	gm.p1.wc <- NewSMStartGame(1, gm.p2.name, gs)
	gm.p2.wc <- NewSMStartGame(2, gm.p1.name, gs)

	for {
		select {
		case msg := <-gm.p1.rc:
			gm.HandleMessage(1, msg)
		case msg := <-gm.p2.rc:
			gm.HandleMessage(2, msg)
		case <-gm.p1.disconnected:
			gm.p2.wc <- NewSMOpponentDisconnected()
			return
		case <-gm.p2.disconnected:
			gm.p1.wc <- NewSMOpponentDisconnected()
			return
		}
	}
}

func (gm *ConcurrentGameManager) HandleMessage(player int, msg ClientMessage) {
	if msg.Msg != ClientMessagePlaceMarker {
		return
	}

	err := gm.g.PlaceMark(player, msg.Placement.Row, msg.Placement.Col)
	if err == nil {
		s := gm.g.State()
		gm.p1.wc <- NewSMUpdateGame(s)
		gm.p2.wc <- NewSMUpdateGame(s)
	}
}
