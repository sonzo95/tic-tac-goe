package server

import (
	"container/list"
	"sync"

	"github.com/gorilla/websocket"
)

type WsMatchmaker struct {
	connQueue list.List // *websocket.Conn
	connLock  sync.Mutex

	gamePool list.List // *ConcurrentGameManager
	// gamePoolLock sync.Mutex
}

func (m *WsMatchmaker) Enqueue(c *websocket.Conn) {
	m.connLock.Lock()
	m.connQueue.PushBack(c)
	for m.connQueue.Len() >= 2 {
		pl1Elem := m.connQueue.Front()
		m.connQueue.Remove(pl1Elem)
		pl2Elem := m.connQueue.Front()
		m.connQueue.Remove(pl2Elem)

		pl1Conn := pl1Elem.Value.(*websocket.Conn)
		pl2Conn := pl2Elem.Value.(*websocket.Conn)
		gm := NewConcurrentGameManager(&WsBroadcaster{
			conns: []*websocket.Conn{pl1Conn, pl2Conn},
			lock:  sync.Mutex{},
		})
		m.gamePool.PushBack(&gm)
	}
	m.connLock.Unlock()
}
