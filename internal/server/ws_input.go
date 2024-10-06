package server

// Object that the server is able to precess when sent by a client via websockets
type InputCommand struct {
	Player int
	Row    int
	Col    int
}
