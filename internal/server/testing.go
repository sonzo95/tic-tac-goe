package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

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
	t.Helper()
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

func assertMessage[T comparable](t testing.TB, conn *websocket.Conn, want T) {
	t.Helper()

	msgChan := readMessage(t, conn)

	select {
	case msg := <-msgChan:
		var got T
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
	case <-time.After(1 * time.Second):
		t.Errorf("timeout exceeded while waiting for message on web socket")
		return
	}
}

func assertConnectionIsClosed(t testing.TB, c *websocket.Conn) {
	t.Helper()
	_, _, err := c.ReadMessage()
	if err == nil {
		t.Errorf("excpected to get an error when reading from client connection")
	}
	if _, ok := err.(*websocket.CloseError); !ok {
		t.Errorf("excpected to close client connection, got error %s", err)
	}
}
