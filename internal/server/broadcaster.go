package server

import (
	"sync"

	"github.com/gorilla/websocket"
	"stefano.sonzogni/tic-tac-toe/internal/game"
)

type Broadcaster interface {
	BroadcastGameState(gs game.GameState)
}

type WsBroadcaster struct {
	conns []*websocket.Conn
	lock  sync.Mutex
}

func (b *WsBroadcaster) BroadcastGameState(gs game.GameState) {
	b.lock.Lock()
	defer b.lock.Unlock()

	for i, conn := range b.conns {
		// TODO: the broadcaster shouldn't be responsible to assign player ids to connections,
		// but for the time being we can get away with this
		conn.WriteJSON(StateUpdate{gs, i + 1})
	}
}
