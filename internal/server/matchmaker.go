package server

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type gameStarter func(p1, p2 *player)

type WsMatchmaker struct {
	playerQueue GenericList[*player]
	connLock    sync.Mutex
	connTimeout time.Duration
	gs          gameStarter
}

func NewWsMatchmaker() *WsMatchmaker {
	return &WsMatchmaker{
		connTimeout: 10 * time.Second,
		gs:          makeGame,
	}
}

func (m *WsMatchmaker) HandleConnection(c *websocket.Conn) {
	slog.Info("handling new connection")
	p := NewPlayer(c)

	select {
	case <-time.After(m.connTimeout):
		slog.Info("connection timed out, didn't receive connect message")
		c.Close()
		return
	case m := <-p.rc:
		if m.Msg == ClientMessageConnect {
			slog.Info("connected player", "name", m.PlayerName)
			p.name = m.PlayerName
		} else {
			slog.Info("received message different from connect, closing connection")
			c.Close()
			return
		}
	}

	c.WriteJSON(NewSMWaitingForMatchmaking())

	m.enqueue(p)
}

func (m *WsMatchmaker) enqueue(p *player) {
	m.connLock.Lock()
	m.playerQueue.PushBack(p)
	for m.playerQueue.Len() >= 2 {
		pl1Conn := m.playerQueue.PopFront()
		pl2Conn := m.playerQueue.PopFront()
		m.gs(pl1Conn, pl2Conn)
	}
	m.connLock.Unlock()
}

func makeGame(pl1, pl2 *player) {
	slog.Info("starting new game", "p1_name", pl1.name, "p2_name", pl2.name)
	gm := NewConcurrentGameManager(pl1, pl2)
	go gm.Start()
}
