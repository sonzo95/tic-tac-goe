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

	t.Run("cannot place markers during enemy's turn", func(t *testing.T) {
		game := NewGame()

		err := game.PlaceMark(CellPlayer2, 0, 0)
		assertError(t, err, ErrInvalidMoveNotPlayerTurn)
	})

	t.Run("cannot place markers on non empty cells", func(t *testing.T) {
		game := NewGame()

		game.PlaceMark(CellPlayer1, 0, 0)
		err := game.PlaceMark(CellPlayer2, 0, 0)
		assertError(t, err, ErrInvalidMovePlaceOnNonEmptyCell)
	})

	t.Run("cannot place markers outside of the board", func(t *testing.T) {
		game := NewGame()

		invalidCellSamples := [][2]int{
			{-1, 0},
			{0, -1},
			{3, 0},
			{0, 3},
		}

		for _, coords := range invalidCellSamples {
			err := game.PlaceMark(CellPlayer1, coords[0], coords[1])
			assertError(t, err, ErrInvalidMoveInvalidCell)
		}
	})

	t.Run("making a winning move updates the game state", func(t *testing.T) {
		tests := []struct {
			game           Game
			nextMovePlayer int
			nextMoveRow    int
			nextMoveCol    int
		}{
			{
				Game{
					GameState{
						CellPlayer1,
						Board{
							{1, 1, 0},
							{0, 0, 0},
							{0, 0, 0},
						},
						WinnerPlayingId,
					},
				},
				CellPlayer1, 0, 2,
			},
			{
				Game{
					GameState{
						CellPlayer1,
						Board{
							{1, 0, 0},
							{1, 0, 0},
							{0, 0, 0},
						},
						WinnerPlayingId,
					},
				},
				CellPlayer1, 2, 0,
			},
			{
				Game{
					GameState{
						CellPlayer1,
						Board{
							{1, 0, 0},
							{0, 0, 0},
							{0, 0, 1},
						},
						WinnerPlayingId,
					},
				},
				CellPlayer1, 1, 1,
			},
			{
				Game{
					GameState{
						CellPlayer2,
						Board{
							{0, 0, 2},
							{0, 2, 0},
							{0, 0, 0},
						},
						WinnerPlayingId,
					},
				},
				CellPlayer2, 2, 0,
			},
		}

		for _, test := range tests {
			err := test.game.PlaceMark(test.nextMovePlayer, test.nextMoveRow, test.nextMoveCol)
			assertError(t, err, nil)

			gotWinner := test.game.State().Winner
			wantWinner := test.nextMovePlayer
			if gotWinner != wantWinner {
				t.Errorf("expected to see %d as winner, got %d", wantWinner, gotWinner)
			}
		}
	})

	t.Run("when the board fills up without winners the game ends in a draw", func(t *testing.T) {
		game := Game{
			GameState{
				CellPlayer1,
				Board{
					{1, 1, 2},
					{2, 2, 1},
					{1, 2, 0},
				},
				WinnerPlayingId,
			},
		}

		err := game.PlaceMark(1, 2, 2)
		assertError(t, err, nil)

		gotWinner := game.State().Winner
		wantWinner := WinnerDrawId
		if gotWinner != wantWinner {
			t.Errorf("expected to see %d as winner, got %d", wantWinner, gotWinner)
		}
	})

	t.Run("no player can make a move after the game ended", func(t *testing.T) {
		game := Game{
			GameState{
				CellPlayer1,
				Board{
					{1, 1, 1},
					{0, 0, 0},
					{0, 0, 0},
				},
				CellPlayer1,
			},
		}

		for _, player := range [2]int{CellPlayer1, CellPlayer2} {
			err := game.PlaceMark(player, 2, 2)
			assertError(t, err, ErrInvalidMoveGameOver)
		}
	})

	t.Run("test draw bug", func(t *testing.T) {
		game := Game{
			GameState{
				CellPlayer2,
				Board{
					{1, 1, 0},
					{0, 2, 0},
					{0, 0, 0},
				},
				WinnerPlayingId,
			},
		}

		err := game.PlaceMark(2, 1, 0)
		assertError(t, err, nil)

		gotWinner := game.State().Winner
		wantWinner := WinnerPlayingId
		if gotWinner != wantWinner {
			t.Errorf("expected to see %d as winner, got %d", wantWinner, gotWinner)
		}
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

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("expected error %v, got %v", want, got)
	}
}
