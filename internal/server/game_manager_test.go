package server

import (
	"slices"
	"testing"
	"time"

	"stefano.sonzogni/tic-tac-toe/internal/game"
)

func TestGameManager(t *testing.T) {
	t.Run("sends startGame messages and handles valid client commands", func(t *testing.T) {
		g := gameSpy{}
		p1 := newPlayer()
		p1.name = "p1"
		p2 := newPlayer()
		p2.name = "p2"
		gm := ConcurrentGameManager{&g, p1, p2}
		go gm.Start()

		time.Sleep(10 * time.Millisecond)

		assertServerMessageOnWriteQueue(t, p1, 1, NewSMStartGame(1, "p2", g.State()))
		assertServerMessageOnWriteQueue(t, p2, 2, NewSMStartGame(2, "p1", g.State()))

		p1.rc <- newPlaceMarkerMessage(0, 0)

		time.Sleep(10 * time.Millisecond)

		p2.rc <- newPlaceMarkerMessage(1, 0)

		time.Sleep(10 * time.Millisecond)

		t.Run("sends inputs to game", func(t *testing.T) {
			want := []move{{1, 0, 0}, {2, 1, 0}}
			assertMoves(t, g, want)
		})

		t.Run("broadcasts updates to players", func(t *testing.T) {
			want := NewSMUpdateGame(game.GameState{
				CurrentPlayer: 2,
				Board:         [3][3]int{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}},
				Winner:        0,
			})
			assertServerMessageOnWriteQueue(t, p1, 1, want)
			assertServerMessageOnWriteQueue(t, p2, 2, want)

			want = NewSMUpdateGame(game.GameState{
				CurrentPlayer: 1,
				Board:         [3][3]int{{1, 0, 0}, {2, 0, 0}, {0, 0, 0}},
				Winner:        0,
			})
			assertServerMessageOnWriteQueue(t, p1, 1, want)
			assertServerMessageOnWriteQueue(t, p2, 2, want)
		})
	})

	t.Run("does not broadcast game updates after invalid moves", func(t *testing.T) {
		g := game.NewGame()
		p1 := newPlayer()
		p2 := newPlayer()
		gm := ConcurrentGameManager{&g, p1, p2}
		go gm.Start()

		time.Sleep(10 * time.Millisecond)
		<-p1.wc
		<-p2.wc

		p1.rc <- newPlaceMarkerMessage(-1, 0)

		time.Sleep(10 * time.Millisecond)

		assertNoServerMessageOnWriteQueue(t, p1, 1)
		assertNoServerMessageOnWriteQueue(t, p2, 2)
	})

	t.Run("ignores unexpected client messages", func(t *testing.T) {
		g := gameSpy{}
		p1 := newPlayer()
		p2 := newPlayer()
		gm := ConcurrentGameManager{&g, p1, p2}
		go gm.Start()

		time.Sleep(10 * time.Millisecond)
		<-p1.wc
		<-p2.wc

		invalidMsg := ClientMessage{
			Msg: ClientMessageConnect,
			Placement: MarkerPlacement{
				Row: 0,
				Col: 0,
			},
		}

		p1.rc <- invalidMsg
		p2.rc <- invalidMsg

		time.Sleep(10 * time.Millisecond)

		want := []move{}
		assertMoves(t, g, want)
	})
}

func assertMoves(t testing.TB, g gameSpy, want []move) {
	t.Helper()
	got := g.Moves
	if !slices.Equal(got, want) {
		t.Errorf("expected game to register %v moves, found %v", want, got)
	}
}

func assertServerMessageOnWriteQueue(t testing.TB, p *player, pid int, want ServerMessage) {
	t.Helper()
	select {
	case got := <-p.wc:
		if want != got {
			t.Errorf("expected to write message %v to player %d, found %v", want, pid, got)
		}
	case <-time.After(time.Second):
		t.Errorf("expected to write message %v to player %d, found none", want, pid)
	}
}

func assertNoServerMessageOnWriteQueue(t testing.TB, p *player, pid int) {
	t.Helper()
	select {
	case got := <-p.wc:
		t.Errorf("expected to not write messages to player %d, found %v", pid, got)
	case <-time.After(time.Second):
		return
	}
}

type move struct {
	player, row, col int
}

type gameSpy struct {
	Moves []move
}

func (g *gameSpy) State() game.GameState {
	gs := initialGameState()
	for _, move := range g.Moves {
		gs.Board[move.row][move.col] = move.player
		if gs.CurrentPlayer == 1 {
			gs.CurrentPlayer = 2
		} else {
			gs.CurrentPlayer = 1
		}
	}
	return gs
}

func (g *gameSpy) PlaceMark(player, row, col int) error {
	g.Moves = append(g.Moves, move{player, row, col})
	return nil
}

func newPlayer() *player {
	return &player{
		rc:           make(chan ClientMessage, 16),
		wc:           make(chan ServerMessage, 16),
		disconnected: false,
	}
}

func newPlaceMarkerMessage(row, col int) ClientMessage {
	return ClientMessage{
		Msg: ClientMessagePlaceMarker,
		Placement: MarkerPlacement{
			Row: row,
			Col: col,
		},
	}
}
