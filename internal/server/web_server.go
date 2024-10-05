package server

import (
	"container/list"
	"sync"

	"github.com/gorilla/websocket"
)

type GenericList[T any] struct {
	data list.List
}

func (l *GenericList[T]) Len() int {
	return l.data.Len()
}

func (l *GenericList[T]) PushBack(t T) {
	l.data.PushBack(t)
}

func (l *GenericList[T]) PopFront() T {
	f := l.data.Front()
	v := f.Value.(T)
	l.data.Remove(f)
	return v
}

type WsMatchmaker struct {
	connQueue GenericList[*websocket.Conn]
	connLock  sync.Mutex

	gamePool GenericList[*ConcurrentGameManager]
	// gamePoolLock sync.Mutex
}

func (m *WsMatchmaker) Enqueue(c *websocket.Conn) {
	m.connLock.Lock()
	m.connQueue.PushBack(c)
	for m.connQueue.Len() >= 2 {
		pl1Conn := m.connQueue.PopFront()
		pl2Conn := m.connQueue.PopFront()

		gm := NewConcurrentGameManager(&WsBroadcaster{
			conns: []*websocket.Conn{pl1Conn, pl2Conn},
			lock:  sync.Mutex{},
		})
		m.gamePool.PushBack(&gm)

		go func() {
			gm.Start()
		}()
	}
	m.connLock.Unlock()
}
