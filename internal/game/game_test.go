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

		if state.CurrentPlayer != CellPlayer1 {
			t.Errorf("expected to be the turn of player 1, but it was %d", state.CurrentPlayer)
		}
		for i, row := range state.Board {
			for j, cellValue := range row {
				if cellValue != CellEmpty {
					t.Errorf("expected cell %d %d to be %d, found %d", i, j, CellEmpty, cellValue)
				}
			}
		}
	})
}
