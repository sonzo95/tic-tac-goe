package server

import (
	"stefano.sonzogni/tic-tac-toe/internal/game"
)

const (
	ClientMessageConnect     = "connect"
	ClientMessagePlaceMarker = "placeMarker"

	ServerMessageWaitingForMatchmaking = "waitingForMatchmaking"
	ServerMessageStartGame             = "startGame"
	ServerMessageUpdateGame            = "updateGame"
	ServerMessageOpponentDisconnected  = "opponentDisconnected"
)

type ClientMessage struct {
	Msg        string          `json:"msg"`
	PlayerName string          `json:"player_name"`
	Placement  MarkerPlacement `json:"marker_placement"`
}

func NewCMConnect(pName string) ClientMessage {
	return ClientMessage{
		Msg:        ClientMessageConnect,
		PlayerName: pName,
	}
}

func NewCMPlaceMarker(row, col int) ClientMessage {
	return ClientMessage{
		Msg: ClientMessageConnect,
		Placement: MarkerPlacement{
			Row: row,
			Col: col,
		},
	}
}

type ServerMessage struct {
	Msg              string         `json:"msg"`
	AssignedPlayerId int            `json:"assigned_player_id"`
	OpponentName     string         `json:"opponent_name"`
	GameState        game.GameState `json:"game_state"`
}

func NewSMWaitingForMatchmaking() ServerMessage {
	return ServerMessage{
		Msg: ServerMessageWaitingForMatchmaking,
	}
}

func NewSMStartGame(pid int, opponentName string, g game.GameState) ServerMessage {
	return ServerMessage{
		Msg:              ServerMessageStartGame,
		AssignedPlayerId: pid,
		OpponentName:     opponentName,
		GameState:        g,
	}
}

func NewSMUpdateGame(g game.GameState) ServerMessage {
	return ServerMessage{
		Msg:       ServerMessageUpdateGame,
		GameState: g,
	}
}

func NewSMOpponentDisconnected() ServerMessage {
	return ServerMessage{
		Msg: ServerMessageOpponentDisconnected,
	}
}

type MarkerPlacement struct {
	Row int
	Col int
}

// Object that the server broadcasts to all connected players
type StateUpdate struct {
	State            game.GameState
	AssignedPlayerId int
}
