package main

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
	"stefano.sonzogni/tic-tac-toe/internal/client"
	"stefano.sonzogni/tic-tac-toe/internal/server"
)

func main() {
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn := connectToServer(5001)
	serverUpdatesCh := make(chan server.ServerMessage, 16)
	serverCommandsCh := make(chan server.ClientMessage, 16)
	go readUpdates(conn, serverUpdatesCh)
	go writeInputs(conn, serverCommandsCh)

	ui := client.UI{
		DefaultFg:     termbox.ColorWhite,
		DefaultBg:     termbox.ColorDefault,
		HighlightedFg: termbox.ColorDefault,
		HighlightedBg: termbox.ColorWhite,
	}

	userCommandCh := make(chan client.Command, 10)
	go client.ListenKeyboard(userCommandCh)

	g := client.NewGame(ui, userCommandCh, serverUpdatesCh, serverCommandsCh)
	g.Start()

	termbox.Close()
}

func connectToServer(port int) *websocket.Conn {
	u := fmt.Sprintf("ws://localhost:%d/play", port)
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)

	if err != nil {
		fmt.Println("Failed to connect to server")
		os.Exit(1)
	}

	return conn
}

func readUpdates(conn *websocket.Conn, out chan server.ServerMessage) {
	for {
		var update server.ServerMessage
		err := conn.ReadJSON(&update)
		if err != nil {
			fmt.Printf("Failed to decode server message: %v", err)
			continue
		}
		out <- update
	}
}

func writeInputs(conn *websocket.Conn, in chan server.ClientMessage) {
	for {
		input := <-in
		err := conn.WriteJSON(input)
		if err != nil {
			fmt.Printf("Failed to write message to server: %v", err)
		}
	}
}
