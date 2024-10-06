package server

import (
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"stefano.sonzogni/tic-tac-toe/internal/game"
)

func TestMatchmaker(t *testing.T) {
	t.Run("every two connections it should create a game and remove them from the queue", func(t *testing.T) {
		ms := WsMatchmaker{}
		server := spinUpServer(t, func(c *websocket.Conn) {
			ms.Enqueue(c)
		})
		defer server.Close()

		n := 50
		conns := []*websocket.Conn{}
		for range n {
			c := dialClient(t, server)
			conns = append(conns, c)
		}

		if ms.connQueue.Len() != 0 {
			t.Errorf("expected matchmaking to empty the whole queue, found %d connections pending", ms.connQueue.Len())
			return
		}

		// assert that all games are also started
		wg := sync.WaitGroup{}
		wg.Add(n)
		for i, conn := range conns {
			msgChan := readMessage(t, conn)
			assertMessage(t, msgChan, StateUpdate{
				State:            initialGameState(),
				AssignedPlayerId: (i % 2) + 1,
			})
			wg.Done()
		}
	})

	t.Run("should propagate ws client messages to games", func(t *testing.T) {
		ms := WsMatchmaker{}
		server := spinUpServer(t, func(c *websocket.Conn) {
			ms.Enqueue(c)
		})
		defer server.Close()

		pl1 := dialClient(t, server)
		pl2 := dialClient(t, server)

		msgChan := readMessage(t, pl1)
		assertMessage(t, msgChan, StateUpdate{initialGameState(), 1})
		msgChan = readMessage(t, pl2)
		assertMessage(t, msgChan, StateUpdate{initialGameState(), 2})

		pl1.WriteJSON(InputCommand{1, 0, 0})

		want := game.GameState{CurrentPlayer: 2, Board: [3][3]int{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}, Winner: 0}
		msgChan = readMessage(t, pl1)
		assertMessage(t, msgChan, StateUpdate{want, 1})
		msgChan = readMessage(t, pl2)
		assertMessage(t, msgChan, StateUpdate{want, 2})
	})
}

func initialGameState() game.GameState {
	g := game.NewGame()
	return g.State()
}
