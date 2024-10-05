package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"stefano.sonzogni/tic-tac-toe/internal/game"
)

func TestBroadcaster(t *testing.T) {
	t.Run("ws broadcaster can register listeners and sends messages to them", func(t *testing.T) {
		b := WsBroadcaster{}
		want := game.GameState{CurrentPlayer: 1, Board: [3][3]int{{0, 0, 0}, {1, 0, 0}, {1, 2, 0}}, Winner: 1}
		nConn := 5

		server := spinUpServer(t, func(c *websocket.Conn) {
			b.AddListener(c)
		})
		defer server.Close()

		conns := []*websocket.Conn{}
		for range nConn {
			conn := dialClient(t, server)
			defer conn.Close()
			conns = append(conns, conn)
		}

		go func() {
			b.BroadcastGameState(want)
		}()

		wg := sync.WaitGroup{}
		wg.Add(nConn)
		for _, conn := range conns {
			msgChan := readMessage(t, conn)
			assertMessage(t, msgChan, want)
			wg.Done()
		}
	})
}

func spinUpServer(t testing.TB, onUpgrade func(*websocket.Conn)) *httptest.Server {
	t.Helper()
	var upgrader = websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("failed to upgrade connection to ws, %v", err)
			return
		}
		onUpgrade(conn)
	}))
	return server
}

func dialClient(t testing.TB, server *httptest.Server) *websocket.Conn {
	t.Helper()
	u := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Errorf("failed to connect to ws server, %v", err)
		return nil
	}
	return conn
}

func readMessage(t testing.TB, conn *websocket.Conn) chan string {
	msgChan := make(chan string, 1)
	go func() {
		_, p, err := conn.ReadMessage()
		if err != nil {
			t.Errorf("failed to read message from ws, %v", err)
			return
		}
		msgChan <- string(p)
	}()
	return msgChan
}

func assertMessage(t testing.TB, msgChan chan string, want game.GameState) {
	select {
	case msg := <-msgChan:
		// todo assert on msg contents
		var got game.GameState
		err := json.NewDecoder(strings.NewReader(msg)).Decode(&got)
		if err != nil {
			t.Errorf("failed to decode ws msg %s, %v", msg, err)
			return
		}

		if got != want {
			t.Errorf("expected message %v, got %v", want, got)
			return
		}
		break
	case <-time.After(2 * time.Second):
		t.Errorf("timeout exceeded while waiting for message on web socket")
		return
	}
}
