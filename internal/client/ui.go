package client

import (
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"stefano.sonzogni/tic-tac-toe/internal/game"
)

type UI struct {
	defaultFg, defaultBg, highlightedFg, highlightedBg termbox.Attribute
}

const (
	topRow    = "┌───┬───┬───┐"
	middleRow = "├───┼───┼───┤"
	bottomRow = "└───┴───┴───┘"
)

func makeBoardRow(row [3]int) string {
	return fmt.Sprintf("│ %d │ %d │ %d │", row[0], row[1], row[2])
}

func (ui *UI) printBoard(x, y int, board game.Board) {
	tbprint(x, y+0, ui.defaultFg, ui.defaultBg, topRow)
	tbprint(x, y+1, ui.defaultFg, ui.defaultBg, makeBoardRow(board[0]))
	tbprint(x, y+2, ui.defaultFg, ui.defaultBg, middleRow)
	tbprint(x, y+3, ui.defaultFg, ui.defaultBg, makeBoardRow(board[1]))
	tbprint(x, y+4, ui.defaultFg, ui.defaultBg, middleRow)
	tbprint(x, y+5, ui.defaultFg, ui.defaultBg, makeBoardRow(board[2]))
	tbprint(x, y+6, ui.defaultFg, ui.defaultBg, bottomRow)
}

func (ui *UI) highlight(x, y int) {
	cell := termbox.GetCell(x, y)
	char := cell.Ch
	tbprint(x, y, ui.highlightedFg, ui.highlightedBg, string(char))
}

// Prints a text on the screen
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}
