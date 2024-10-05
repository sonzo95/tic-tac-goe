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
}

// TODO: do i actually need to save all games in the array?
func (m *WsMatchmaker) Enqueue(c *websocket.Conn) {
	m.connLock.Lock()
	m.connQueue.PushBack(c)
	for m.connQueue.Len() >= 2 {
		pl1Conn := m.connQueue.PopFront()
		pl2Conn := m.connQueue.PopFront()
		makeGame(pl1Conn, pl2Conn)
	}
	m.connLock.Unlock()
}

func makeGame(pl1, pl2 *websocket.Conn) {
	gm := NewConcurrentGameManager(&WsBroadcaster{
		conns: []*websocket.Conn{pl1, pl2},
		lock:  sync.Mutex{},
	})

	gm.Start()

	listenInputs(pl1, &gm)
	listenInputs(pl2, &gm)
}

type InputCommand struct {
	Player int
	Row    int
	Col    int
}

func listenInputs(c *websocket.Conn, gm GameManager) {
	go func() {
		// TODO: check if connection is still open?
		for {
			var input InputCommand
			err := c.ReadJSON(&input)
			if err == nil {
				gm.HandleMessage(input.Player, input.Row, input.Col)
			}
		}
	}()
}
