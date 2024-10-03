package game

import "errors"

/*
Tic Tac Toe

3x3 field
2 players
turns
*/

// Cells are indexed by [row][col], with indices gronwind downwards and rightwards
type Board [3][3]int

func emptyBoard() Board {
	return Board{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
}

const (
	CellEmpty   = 0
	CellPlayer1 = 1
	CellPlayer2 = 2
)

var (
	ErrInvalidMoveNotPlayerTurn       = errors.New("invalid move: not player turn")
	ErrInvalidMovePlaceOnNonEmptyCell = errors.New("invalid move: placing marker on non empty cell")
	ErrInvalidMoveInvalidCell         = errors.New("invalid move: invalid cell coordinates")
)

type GameState struct {
	CurrentPlayer int
	Board         [3][3]int
}

type Game struct {
	state GameState
}

func NewGame() Game {
	return Game{
		state: GameState{
			CurrentPlayer: CellPlayer1,
			Board:         emptyBoard(),
		},
	}
}

// Returns a copy of the game state
func (g *Game) State() GameState {
	return g.state
}

func (g *Game) PlaceMark(player, row, col int) error {
	if g.state.CurrentPlayer != player {
		return ErrInvalidMoveNotPlayerTurn
	}

	if row < 0 || row > 2 || col < 0 || col > 2 {
		return ErrInvalidMoveInvalidCell
	}

	if g.state.Board[row][col] != CellEmpty {
		return ErrInvalidMovePlaceOnNonEmptyCell
	}

	g.state.Board[row][col] = player
	if g.state.CurrentPlayer == CellPlayer1 {
		g.state.CurrentPlayer = CellPlayer2
	} else {
		g.state.CurrentPlayer = CellPlayer1
	}
	return nil
}
