package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"stefano.sonzogni/tic-tac-toe/internal/server"
)

func main() {
	matchmaker := server.NewWsMatchmaker()

	mux := http.NewServeMux()

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	mux.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("received play request")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Info("connection established")
		matchmaker.HandleConnection(conn)
	})

	err := http.ListenAndServe(":5001", mux)
	if err != nil {
		slog.Error("failed to start server", "err", err.Error())
	}
}
