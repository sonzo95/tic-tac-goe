package client

import (
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"stefano.sonzogni/tic-tac-toe/internal/game"
)

type UI struct {
	DefaultFg, DefaultBg, HighlightedFg, HighlightedBg termbox.Attribute
}

type GameRenderer interface {
	RenderGame(s game.GameState, cell Cell, msg string, playerId int)
}

const (
	topRow    = "┌───┬───┬───┐"
	middleRow = "├───┼───┼───┤"
	bottomRow = "└───┴───┴───┘"
)

func makeBoardRow(row [3]int) string {
	return fmt.Sprintf("│ %d │ %d │ %d │", row[0], row[1], row[2])
}

func (ui UI) RenderGame(s game.GameState, cell Cell, msg string, playerId int) {
	termbox.Clear(ui.DefaultFg, ui.DefaultBg)

	if playerId == s.CurrentPlayer {
		tbprint(0, 0, ui.DefaultFg, ui.DefaultBg, "It's your turn!")
	} else {
		tbprint(0, 0, ui.DefaultFg, ui.DefaultBg, "Waiting for the opponent to move")
	}

	ui.printBoard(0, 2, s.Board)
	x, y := cell.toScreenCoords()
	ui.highlight(x, y+2)

	switch s.Winner {
	case game.WinnerPlayingId:
		tbprint(0, 10, ui.DefaultFg, ui.DefaultBg, msg)
	case game.WinnerDrawId:
		tbprint(0, 10, ui.DefaultFg, ui.DefaultBg, "It's a draw!")
	case playerId:
		tbprint(0, 10, ui.DefaultFg, ui.DefaultBg, "You won!")
	default:
		tbprint(0, 10, ui.DefaultFg, ui.DefaultBg, "You lost!")
	}

	termbox.Flush()
}

func (ui *UI) printBoard(x, y int, board game.Board) {
	tbprint(x, y+0, ui.DefaultFg, ui.DefaultBg, topRow)
	tbprint(x, y+1, ui.DefaultFg, ui.DefaultBg, makeBoardRow(board[0]))
	tbprint(x, y+2, ui.DefaultFg, ui.DefaultBg, middleRow)
	tbprint(x, y+3, ui.DefaultFg, ui.DefaultBg, makeBoardRow(board[1]))
	tbprint(x, y+4, ui.DefaultFg, ui.DefaultBg, middleRow)
	tbprint(x, y+5, ui.DefaultFg, ui.DefaultBg, makeBoardRow(board[2]))
	tbprint(x, y+6, ui.DefaultFg, ui.DefaultBg, bottomRow)
}

func (ui *UI) highlight(x, y int) {
	cell := termbox.GetCell(x, y)
	char := cell.Ch
	tbprint(x, y, ui.HighlightedFg, ui.HighlightedBg, string(char))
}

// Prints a text on the screen
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

// Cell helper

func (c Cell) toScreenCoords() (int, int) {
	// 0 1 2 -> 2 6 10
	x := 2 + c.r*4
	// 0 1 2 -> 1 3 5
	y := c.c*2 + 1
	return x, y
}
