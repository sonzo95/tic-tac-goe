package game

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

func (g *Game) PlaceMark(player, row, col int) {

}
