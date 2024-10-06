package main

import (
	"fmt"
	"os"

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

	ui := client.UI{
		DefaultFg:     termbox.ColorWhite,
		DefaultBg:     termbox.ColorDefault,
		HighlightedFg: termbox.ColorDefault,
		HighlightedBg: termbox.ColorWhite,
	}

	commandCh := make(chan client.Command, 10)
	client.ListenKeyboard(commandCh)

	g := client.NewGame(ui, commandCh, make(chan server.StateUpdate))
	g.Start()

	termbox.Close()
}
