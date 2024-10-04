package server

import (
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

func TestGameManager(t *testing.T) {
	t.Run("updates game state if commands are received", func(t *testing.T) {
		gs := gameSpy{}
		gm := ConcurrentGameManager{&gs, sync.Mutex{}}

		player := 1
		row := 0
		col := 0

		gm.HandleMessage(player, row, col)

		if len(gs.Moves) != 1 {
			t.Errorf("expected to record one move, found %d", len(gs.Moves))
		}
	})

	t.Run("updates game state safely wrt concurrency", func(t *testing.T) {
		gs := gameSpy{}
		gm := ConcurrentGameManager{&gs, sync.Mutex{}}

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

		if len(gs.Moves) != numOfMoves {
			t.Errorf("expected to record %d moves, found %d", numOfMoves, len(gs.Moves))
		}
	})
}
