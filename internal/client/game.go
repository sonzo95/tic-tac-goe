package client

import (
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
	"stefano.sonzogni/tic-tac-toe/internal/game"
)

type Game struct {
	// state, current highlighted cell, ui, channels?
}

func (g *Game) Start() {
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ui := UI{termbox.ColorWhite, termbox.ColorDefault, termbox.ColorDefault, termbox.ColorWhite}
	ui.printBoard(0, 0, game.Board{{1, 2, 0}, {1, 2, 2}, {0, 0, 0}})
	ui.highlight(2, 1)
	termbox.Flush()

	time.Sleep(4 * time.Second)
	termbox.Close()
}
