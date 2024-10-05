package server

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestMatchmaker(t *testing.T) {
	t.Run("every two connections it should create a game and remove them from the queue", func(t *testing.T) {
		ms := WsMatchmaker{}
		server := spinUpServer(t, func(c *websocket.Conn) {
			ms.Enqueue(c)
		})
		defer server.Close()

		n := 50
		for range n {
			go func() {
				dialClient(t, server)
			}()
		}

		time.Sleep(50 * time.Millisecond)

		if ms.connQueue.Len() != 0 {
			t.Errorf("expected matchmaking to empty the whole queue, found %d connections pending", ms.connQueue.Len())
			return
		}
		if ms.gamePool.Len() != n/2 {
			t.Errorf("expected matchmaking to create %d games, found %d games", n/2, ms.gamePool.Len())
			return
		}
	})
}
