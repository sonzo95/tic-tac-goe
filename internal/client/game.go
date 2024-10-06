package client

import (
	"stefano.sonzogni/tic-tac-toe/internal/game"
	"stefano.sonzogni/tic-tac-toe/internal/server"
)

type Game struct {
	ui               GameRenderer
	commandCh        chan Command
	updatesCh        chan server.StateUpdate
	cursorX, cursorY int
	msg              string
	quit             bool
	playerId         int
	state            game.GameState
}

func NewGame(ui GameRenderer, commandCh chan Command, updatesCh chan server.StateUpdate) *Game {
	return &Game{
		ui:        ui,
		commandCh: commandCh,
		updatesCh: updatesCh,
	}
}

func (g *Game) Start() {
	g.render()

	for !g.quit {
		g.msg = ""

		// should become select on commands and server events
		select {
		case cmd := <-g.commandCh:
			cmd(g)
		case update := <-g.updatesCh:
			g.playerId = update.AssignedPlayerId
			g.state = update.State
		}

		g.render()
	}
}

func (g *Game) render() {
	g.ui.RenderGame(
		g.state,
		Cell{g.cursorX, g.cursorY},
		g.msg,
		g.playerId,
	)
}

type Cell struct {
	r, c int
}
