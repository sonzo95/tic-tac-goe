package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"stefano.sonzogni/tic-tac-toe/internal/game"
)

func TestHandleConnection(t *testing.T) {
	t.Run("should close the connection if no connect message is received within timeout duration", func(t *testing.T) {
		ms := WsMatchmaker{
			connTimeout: time.Millisecond * 10,
		}
		server := spinUpServer(t, func(c *websocket.Conn) {
			ms.HandleConnection(c)
		})
		defer server.Close()

		c := dialClient(t, server)

		time.Sleep(time.Millisecond * 15)

		assertConnectionIsClosed(t, c)
	})

	t.Run("should close the connection if the wrong message is received", func(t *testing.T) {
		ms := NewWsMatchmaker()
		server := spinUpServer(t, func(c *websocket.Conn) {
			ms.HandleConnection(c)
		})
		defer server.Close()

		c := dialClient(t, server)
		c.WriteJSON(ClientMessage{Msg: ClientMessagePlaceMarker})

		assertConnectionIsClosed(t, c)
	})

	t.Run("should handle clients that sends a connect message", func(t *testing.T) {
		ms := NewWsMatchmaker()
		server := spinUpServer(t, func(c *websocket.Conn) {
			ms.HandleConnection(c)
		})
		defer server.Close()

		c := dialClient(t, server)
		c.WriteJSON(ClientMessage{
			Msg:        ClientMessageConnect,
			PlayerName: "name",
		})

		t.Run("should write matchmaking message", func(t *testing.T) {
			assertMessage(t, c, NewSMWaitingForMatchmaking())
		})

		t.Run("should enqueue the player", func(t *testing.T) {
			if ms.playerQueue.Len() != 1 {
				t.Errorf("expected to find %d players in queue, found %d", 1, ms.playerQueue.Len())
			}
			player := ms.playerQueue.PopFront()
			if player.name != "name" {
				t.Errorf("expected to find %s in queue, found %s", "name", player.name)
			}
		})
	})

	t.Run("every two connections it should create a game and remove them from the queue", func(t *testing.T) {
		ms := NewWsMatchmaker()
		gameStartedCh := make(chan struct{}, 50)
		ms.gs = func(p1, p2 *player) { gameStartedCh <- struct{}{} }

		server := spinUpServer(t, func(c *websocket.Conn) {
			ms.HandleConnection(c)
		})
		defer server.Close()

		n := 50
		for i := range n {
			c := dialClient(t, server)
			c.WriteJSON(ClientMessage{
				Msg:        ClientMessageConnect,
				PlayerName: fmt.Sprintf("p%d", i),
			})
		}

		time.Sleep(10 * time.Millisecond)

		if ms.playerQueue.Len() != 0 {
			t.Errorf("expected matchmaking to empty the whole queue, found %d connections pending", ms.playerQueue.Len())
			return
		}

		// assert that all games are also started
		got := len(gameStartedCh)
		want := n / 2
		if got != want {
			t.Errorf("expected to start %d games, found %d", want, got)
		}
	})
}

func initialGameState() game.GameState {
	g := game.NewGame()
	return g.State()
}
