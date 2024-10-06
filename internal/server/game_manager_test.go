package server

import (
	"slices"
	"sync"
	"testing"
	"time"

	"stefano.sonzogni/tic-tac-toe/internal/game"
)

type move struct {
	player, row, col int
}

type gameSpy struct {
	Moves []move
}

func (g *gameSpy) State() game.GameState {
	return game.GameState{}
}

func (g *gameSpy) PlaceMark(player, row, col int) error {
	g.Moves = append(g.Moves, move{player, row, col})
	return nil
}

type broadcasterSpy struct {
	gameStates []game.GameState
}

func (b *broadcasterSpy) BroadcastGameState(gs game.GameState) {
	b.gameStates = append(b.gameStates, gs)
}

func TestGameManager(t *testing.T) {
	t.Run("updates game state if commands are received", func(t *testing.T) {
		g := gameSpy{}
		gm := ConcurrentGameManager{&g, sync.Mutex{}, &broadcasterSpy{}}

		player := 1
		row := 0
		col := 0

		gm.HandleMessage(player, row, col)

		assertNumberOfMoves(t, g, 1)
	})

	t.Run("updates game state safely wrt concurrency", func(t *testing.T) {
		g := gameSpy{}
		gm := ConcurrentGameManager{&g, sync.Mutex{}, &broadcasterSpy{}}

		player := 1
		row := 0
		col := 0
		numOfMoves := 100000

		for range numOfMoves {
			go func() {
				gm.HandleMessage(player, row, col)
			}()
		}

		time.Sleep(50 * time.Millisecond)

		assertNumberOfMoves(t, g, numOfMoves)
	})

	t.Run("updates to game state are broadcasted to listeners", func(t *testing.T) {
		g := game.NewGame()
		broadcaster := &broadcasterSpy{}
		gm := ConcurrentGameManager{&g, sync.Mutex{}, broadcaster}

		gm.HandleMessage(1, 0, 0)
		want := []game.GameState{gm.g.State()}
		assertBroadcasts(t, broadcaster, want)

		gm.HandleMessage(2, 1, 0)
		want = append(want, gm.g.State())
		assertBroadcasts(t, broadcaster, want)
	})

	t.Run("invalid moves are ignored and not broadcasted to listeners", func(t *testing.T) {
		g := game.NewGame()
		broadcaster := &broadcasterSpy{}
		gm := ConcurrentGameManager{&g, sync.Mutex{}, broadcaster}

		gm.HandleMessage(2, 0, 0)
		want := []game.GameState{}
		assertBroadcasts(t, broadcaster, want)
	})
}

func assertNumberOfMoves(t testing.TB, gs gameSpy, want int) {
	t.Helper()
	got := len(gs.Moves)
	if got != want {
		t.Errorf("expected to record %d moves, found %d", want, got)
	}
}

func assertBroadcasts(t testing.TB, broadcaster *broadcasterSpy, want []game.GameState) {
	t.Helper()
	if !slices.Equal(broadcaster.gameStates, want) {
		t.Errorf("expected to broadcast %v, found %v", want, broadcaster.gameStates)
	}
}
