package client

import (
	"stefano.sonzogni/tic-tac-toe/internal/game"
	"stefano.sonzogni/tic-tac-toe/internal/server"
)

type Game struct {
	ui GameRenderer
	// input channel that allows to send commands to the game
	userCommandsCh chan Command
	// input channel that allows to send state updates to the game
	serverUpdatesCh chan server.ServerMessage
	// output channel that allows to send commands to the server
	serverCommandsCh     chan server.ClientMessage
	cursorX, cursorY     int
	msg                  string
	quit                 bool
	playerId             int
	state                game.GameState
	playerName           string
	opponentName         string
	opponentDisconnected bool
}

func NewGame(
	playerId int,
	playerName string,
	opponentName string,
	initialState game.GameState,
	ui GameRenderer,
	userCommandCh chan Command,
	serverUpdatesCh chan server.ServerMessage,
	serverCommandsCh chan server.ClientMessage,
) *Game {
	return &Game{
		ui:               ui,
		userCommandsCh:   userCommandCh,
		serverUpdatesCh:  serverUpdatesCh,
		serverCommandsCh: serverCommandsCh,
		playerId:         playerId,
		playerName:       playerName,
		opponentName:     opponentName,
		state:            initialState,
	}
}

func (g *Game) Start() {
	g.render()

	for !g.quit && !g.opponentDisconnected {
		g.msg = ""

		// should become select on commands and server events
		select {
		case cmd := <-g.userCommandsCh:
			cmd(g)
		case update := <-g.serverUpdatesCh:
			switch update.Msg {
			case server.ServerMessageUpdateGame:
				g.state = update.GameState
			case server.ServerMessageOpponentDisconnected:
				g.opponentDisconnected = true
				g.msg = "The opponent left!"
			}
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
