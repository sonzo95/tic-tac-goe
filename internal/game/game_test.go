package game

import "testing"

/*
1. starting a game gives an empty field
2. starting a game gives turn to player X
3. the current player can make a move and place a symbol on an empty cell
4. making a move shall pass the turn to the other contestant
5. determine winning states
*/

func TestGame(t *testing.T) {
	t.Run("game initialization", func(t *testing.T) {
		game := NewGame()

		var state GameState = game.State()

		assertCurrentPlayer(t, state, CellPlayer1)
		assertBoard(t, state, emptyBoard())
	})

	t.Run("making a valid move places a marker on the board and changes current player", func(t *testing.T) {
		game := NewGame()

		game.PlaceMark(CellPlayer1, 0, 0)
		state := game.State()
		assertBoard(t, state, Board{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}})
		assertCurrentPlayer(t, state, CellPlayer2)

		game.PlaceMark(CellPlayer2, 0, 2)
		state = game.State()
		assertBoard(t, state, Board{{1, 0, 2}, {0, 0, 0}, {0, 0, 0}})
		assertCurrentPlayer(t, state, CellPlayer1)
	})
}

func assertCurrentPlayer(t testing.TB, state GameState, want int) {
	t.Helper()
	if state.CurrentPlayer != want {
		t.Errorf("expected to be the turn of player %d, but it was %d", want, state.CurrentPlayer)
	}
}

func assertBoard(t testing.TB, state GameState, want Board) {
	t.Helper()
	for i := range 3 {
		for j := range 3 {
			gotCell := state.Board[i][j]
			wantCell := want[i][j]
			if gotCell != wantCell {
				t.Errorf("expected cell [%d,%d] to be %d, found %d", i, j, wantCell, gotCell)
			}
		}
	}
}
