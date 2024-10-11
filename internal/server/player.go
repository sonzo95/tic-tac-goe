package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type player struct {
	name         string
	rc           chan ClientMessage
	wc           chan ServerMessage
	disconnected bool
}

func NewPlayer(c *websocket.Conn) *player {
	p := player{
		rc: make(chan ClientMessage, 16),
		wc: make(chan ServerMessage, 16),
	}
	go p.readMessages(c)
	go p.writeMessages(c)
	return &p
}

func (pl *player) readMessages(c *websocket.Conn) {
	for {
		var m ClientMessage
		t, p, e := c.ReadMessage()
		if e != nil {
			c.Close()
			close(pl.rc)
			pl.disconnected = true
			return
		}
		if t == websocket.TextMessage {
			e = json.Unmarshal(p, &m)
			if e == nil {
				pl.rc <- m
			}
		}
	}
}

func (p *player) writeMessages(c *websocket.Conn) {
	for m := range p.wc {
		jm, e := json.Marshal(m)
		if e == nil {
			c.WriteMessage(websocket.TextMessage, jm)
		}
	}
}
