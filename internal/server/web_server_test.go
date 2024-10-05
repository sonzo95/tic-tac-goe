package server

import (
	"sync"
	"testing"
	"time"

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
		connsLock := sync.Mutex{}
		for range n {
			go func() {
				c := dialClient(t, server)
				connsLock.Lock()
				conns = append(conns, c)
				connsLock.Unlock()
			}()
		}

		time.Sleep(50 * time.Millisecond)

		if ms.connQueue.Len() != 0 {
			t.Errorf("expected matchmaking to empty the whole queue, found %d connections pending", ms.connQueue.Len())
			return
		}

		// assert that all games are also started
		wg := sync.WaitGroup{}
		wg.Add(n)
		for _, conn := range conns {
			msgChan := readMessage(t, conn)
			assertMessage(t, msgChan, game.GameState{})
			wg.Done()
		}
	})
}
