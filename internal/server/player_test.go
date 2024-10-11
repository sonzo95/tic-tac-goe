package server

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"stefano.sonzogni/tic-tac-toe/internal/game"
)

func TestPlayer(t *testing.T) {
	serverConn := make(chan *websocket.Conn, 1)
	s := spinUpServer(t, func(conn *websocket.Conn) {
		serverConn <- conn
	})
	defer s.Close()
	clientConn := dialClient(t, s)

	c := <-serverConn
	p := NewPlayer(c)
	p.name = "name"

	t.Run("should write server messages", func(t *testing.T) {
		want := ServerMessage{
			Msg:              "hello",
			AssignedPlayerId: 1,
			OpponentName:     "opp",
			GameState:        game.GameState{},
		}
		p.wc <- want
		msg := <-readMessage(t, clientConn)
		var got ServerMessage
		json.NewDecoder(strings.NewReader(msg)).Decode(&got)

		if want != got {
			t.Errorf("expected to read %v from connection, got %v", want, got)
		}
	})

	t.Run("should read server messages", func(t *testing.T) {
		want := ClientMessage{
			Msg: "hello",
		}
		clientConn.WriteJSON(want)
		got := <-p.rc
		if want != got {
			t.Errorf("expected to read %v from connection, got %v", want, got)
		}
	})

	t.Run("should send disconnect signal if client disconnects", func(t *testing.T) {
		clientConn.Close()

		time.Sleep(10 * time.Millisecond)

		if !p.disconnected {
			t.Errorf("expected client to be disconnected")
		}
	})
}
