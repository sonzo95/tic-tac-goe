package client

import (
	"testing"
	"time"

	"stefano.sonzogni/tic-tac-toe/internal/game"
	"stefano.sonzogni/tic-tac-toe/internal/server"
)

func TestGame(t *testing.T) {
	t.Run("commands get processed and trigger rerender", func(t *testing.T) {
		r := GameRendererSpy{}
		cc := make(chan Command, 1)
		g := NewGame(1, "", "", initialState(), &r, cc, make(chan server.ServerMessage), make(chan server.ClientMessage))
		go g.Start()

		executed := false
		cc <- func(g *Game) {
			executed = true
		}

		time.Sleep(10 * time.Millisecond)
		assertRenderCount(t, r, 2)
		if !executed {
			t.Errorf("expected command to be executed")
		}
	})

	t.Run("state updates trigger rerender", func(t *testing.T) {
		r := GameRendererSpy{}
		su := make(chan server.ServerMessage, 1)
		g := NewGame(2, "", "", initialState(), &r, make(chan Command), su, make(chan server.ClientMessage))
		go g.Start()

		newBoard := game.Board{{1, 1, 1}, {2, 2, 2}, {1, 1, 1}}
		newState := game.GameState{
			Board:         newBoard,
			CurrentPlayer: 2,
			Winner:        2,
		}
		su <- server.NewSMUpdateGame(newState)

		time.Sleep(10 * time.Millisecond)

		assertRenderCount(t, r, 2)
	})

	t.Run("state updates trigger rerender", func(t *testing.T) {
		r := GameRendererSpy{}
		su := make(chan server.ServerMessage, 1)
		g := NewGame(2, "", "", initialState(), &r, make(chan Command), su, make(chan server.ClientMessage))
		go g.Start()

		newBoard := game.Board{{1, 1, 1}, {2, 2, 2}, {1, 1, 1}}
		newState := game.GameState{
			Board:         newBoard,
			CurrentPlayer: 2,
			Winner:        2,
		}
		su <- server.NewSMUpdateGame(newState)

		time.Sleep(10 * time.Millisecond)

		assertRenderCount(t, r, 2)
		assertLastRender(t, r, newState, 2)

		// also ignores state updates from non-update massegaes
		messages := []server.ServerMessage{
			server.NewSMStartGame(1, "", game.GameState{}),
			server.NewSMWaitingForMatchmaking(),
			server.NewSMOpponentDisconnected(),
		}
		for i, msg := range messages {
			su <- msg

			time.Sleep(10 * time.Millisecond)

			assertRenderCount(t, r, 3+i)
			assertLastRender(t, r, newState, 2)
		}
	})

	t.Run("disconnection messages trigger rerender", func(t *testing.T) {
		r := GameRendererSpy{}
		su := make(chan server.ServerMessage, 1)
		g := NewGame(2, "", "", initialState(), &r, make(chan Command), su, make(chan server.ClientMessage))
		go g.Start()

		su <- server.NewSMOpponentDisconnected()

		time.Sleep(10 * time.Millisecond)

		assertRenderCount(t, r, 2)
		if !g.opponentDisconnected {
			t.Error("expected opponent to be marked as disconnected")
		}
	})
}

func assertRenderCount(t testing.TB, r GameRendererSpy, want int) {
	t.Helper()
	got := len(r.renders)
	if got != want {
		t.Errorf("expected %d renders, found %d", want, got)
	}
}

func assertLastRender(t testing.TB, r GameRendererSpy, state game.GameState, player int) {
	t.Helper()
	lastRender := r.renders[len(r.renders)-1]
	if lastRender.s != state {
		t.Errorf("expected to render state %v, rendered %v", state, lastRender.s)
	}
	if lastRender.playerId != player {
		t.Errorf("expected to render player %d, rendered %d", player, lastRender.playerId)
	}
}

type RenderData struct {
	s        game.GameState
	cell     Cell
	msg      string
	playerId int
}

type GameRendererSpy struct {
	renders []RenderData
}

func (gr *GameRendererSpy) RenderGame(s game.GameState, cell Cell, msg string, playerId int) {
	gr.renders = append(gr.renders, RenderData{s, cell, msg, playerId})
}

func initialState() game.GameState {
	return game.GameState{
		CurrentPlayer: 1,
		Board:         [3][3]int{},
		Winner:        0,
	}
}
