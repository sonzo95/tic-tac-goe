package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WsMatchmaker struct {
	connQueue GenericList[*websocket.Conn]
	connLock  sync.Mutex
}

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
