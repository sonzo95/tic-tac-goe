package server

import (
	"container/list"
	"sync"

	"github.com/gorilla/websocket"
)

type WsMatchmaker struct {
	connQueue list.List
	connLock  sync.Mutex

	gamePool     list.List
	gamePoolLock sync.Mutex
}

func (m *WsMatchmaker) Enqueue(c *websocket.Conn) {
	m.connLock.Lock()
	m.connQueue.PushBack(c)
	for m.connQueue.Len() >= 2 {
		pl1Elem := m.connQueue.Front()
		m.connQueue.Remove(pl1Elem)
		pl2Elem := m.connQueue.Front()
		m.connQueue.Remove(pl2Elem)

		b := WsBroadcaster{}
		pl1Conn := pl1Elem.Value.(*websocket.Conn)
		pl2Conn := pl2Elem.Value.(*websocket.Conn)
		b.AddListener(pl1Conn)
		b.AddListener(pl2Conn)
		m.gamePool.PushBack(NewConcurrentGameManager(&b))
	}
	m.connLock.Unlock()
}
