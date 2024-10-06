package server

import "stefano.sonzogni/tic-tac-toe/internal/game"

// Object that the server is able to precess when sent by a client via websockets
type InputCommand struct {
	Player int
	Row    int
	Col    int
}

// Object that the server broadcasts to all connected players
type StateUpdate struct {
	State            game.GameState
	AssignedPlayerId int
}
