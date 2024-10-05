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

	for _, conn := range b.conns {
		conn.WriteJSON(gs)
	}
}

func (b *WsBroadcaster) AddListener(conn *websocket.Conn) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.conns = append(b.conns, conn)
}
