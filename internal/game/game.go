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

func checkRowWinningCon(b Board, row int) bool {
	return b[row][0] != CellEmpty &&
		b[row][0] == b[row][1] &&
		b[row][0] == b[row][2]
}

func checkColWinningCon(b Board, col int) bool {
	return b[0][col] != CellEmpty &&
		b[0][col] == b[1][col] &&
		b[0][col] == b[2][col]
}

func checkMainDiagonalWinningCon(b Board, row, col int) bool {
	return row == col &&
		b[0][0] != CellEmpty &&
		b[0][0] == b[1][1] &&
		b[0][0] == b[2][2]
}

func checkSecDiagonalWinningCon(b Board, row, col int) bool {
	return row+col == 2 &&
		b[0][2] != CellEmpty &&
		b[0][2] == b[1][1] &&
		b[0][2] == b[2][0]
}

func isBoardFilled(b Board) bool {
	for row := range 3 {
		for col := range 3 {
			if b[row][col] == CellEmpty {
				return false
			}
		}
	}
	return true
}

const (
	CellEmpty   = 0
	CellPlayer1 = 1
	CellPlayer2 = 2

	// the game is still being played
	WinnerPlayingId = 0
	// the game ended in a draw
	WinnerDrawId = 3
)

var (
	ErrInvalidMoveNotPlayerTurn       = errors.New("invalid move: not player turn")
	ErrInvalidMovePlaceOnNonEmptyCell = errors.New("invalid move: placing marker on non empty cell")
	ErrInvalidMoveInvalidCell         = errors.New("invalid move: invalid cell coordinates")
	ErrInvalidMoveGameOver            = errors.New("invalid move: game is over")
)

type GameState struct {
	CurrentPlayer int
	Board         [3][3]int
	Winner        int
}

func (gs *GameState) swapPlayerTurn() {
	if gs.CurrentPlayer == CellPlayer1 {
		gs.CurrentPlayer = CellPlayer2
	} else {
		gs.CurrentPlayer = CellPlayer1
	}
}

type StatefulInteractableGame interface {
	State() GameState
	PlaceMark(player, row, col int) error
}

type Game struct {
	state GameState
}

func NewGame() Game {
	return Game{
		state: GameState{
			CurrentPlayer: CellPlayer1,
			Board:         emptyBoard(),
			Winner:        WinnerPlayingId,
		},
	}
}

// Returns a copy of the game state
func (g *Game) State() GameState {
	return g.state
}

func (g *Game) PlaceMark(player, row, col int) error {
	if g.state.Winner != WinnerPlayingId {
		return ErrInvalidMoveGameOver
	}

	if g.state.CurrentPlayer != player {
		return ErrInvalidMoveNotPlayerTurn
	}

	if row < 0 || row > 2 || col < 0 || col > 2 {
		return ErrInvalidMoveInvalidCell
	}

	if g.state.Board[row][col] != CellEmpty {
		return ErrInvalidMovePlaceOnNonEmptyCell
	}

	// place marker
	g.state.Board[row][col] = player

	// check winning con
	if checkRowWinningCon(g.state.Board, row) ||
		checkColWinningCon(g.state.Board, col) ||
		checkMainDiagonalWinningCon(g.state.Board, row, col) ||
		checkSecDiagonalWinningCon(g.state.Board, row, col) {
		g.state.Winner = player
		return nil
	}

	// check draw
	if isBoardFilled(g.state.Board) {
		g.state.Winner = WinnerDrawId
		return nil
	}

	g.state.swapPlayerTurn()
	return nil
}
