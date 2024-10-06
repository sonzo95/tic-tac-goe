package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"stefano.sonzogni/tic-tac-toe/internal/game"
)

func main() {
	out := os.Stdout
	inBuffer := bufio.NewScanner(os.Stdin)

	g := game.NewGame()

	for g.State().Winner == game.WinnerPlayingId {
		printCurrentState(&g, out)
		fmt.Fprintf(out, "Insert placement position: ")
		inBuffer.Scan()
		inputStr := strings.Trim(inBuffer.Text(), " ")
		tokens := strings.Split(inputStr, " ")
		if len(tokens) != 2 {
			fmt.Fprintln(out, "Invalid input, please type 'row col'")
			continue
		}

		row, err := strconv.Atoi(tokens[0])
		if err != nil {
			fmt.Fprintln(out, "Invalid input, row must be a number")
			continue
		}

		col, err := strconv.Atoi(tokens[1])
		if err != nil {
			fmt.Fprintln(out, "Invalid input, col must be a number")
			continue
		}

		err = g.PlaceMark(g.State().CurrentPlayer, row, col)
		if err != nil {
			fmt.Fprintf(out, "Invalid move: %v\n", err)
			continue
		}
	}

	printCurrentState(&g, out)
	winner := g.State().Winner
	switch winner {
	case game.CellPlayer1, game.CellPlayer2:
		fmt.Fprintf(out, "Player %d won!\n", winner)
	case game.WinnerDrawId:
		fmt.Fprintln(out, "Game ended in a draw")
	}
}

func printCurrentState(g *game.Game, out io.Writer) {
	fmt.Fprintln(out, "-------")
	for _, row := range g.State().Board {
		rowBuilder := strings.Builder{}
		rowBuilder.WriteString("|")

		for _, cell := range row {
			switch cell {
			case game.CellEmpty:
				rowBuilder.WriteString(" ")
			case game.CellPlayer1:
				rowBuilder.WriteString("X")
			case game.CellPlayer2:
				rowBuilder.WriteString("O")
			default:
				rowBuilder.WriteString("?")
			}
			rowBuilder.WriteString("|")
		}
		fmt.Fprintln(out, rowBuilder.String())
		fmt.Fprintln(out, "-------")
	}
	fmt.Fprintf(out, "It's player %d turn\n", g.State().CurrentPlayer)
}
