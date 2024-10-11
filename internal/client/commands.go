package client

import "stefano.sonzogni/tic-tac-toe/internal/server"

type Command func(g *Game)

func commandLeft(g *Game) {
	if g.cursorX > 0 {
		g.cursorX--
	}
}

func commandRight(g *Game) {
	if g.cursorX < 2 {
		g.cursorX++
	}
}

func commandUp(g *Game) {
	if g.cursorY > 0 {
		g.cursorY--
	}
}

func commandDown(g *Game) {
	if g.cursorY < 2 {
		g.cursorY++
	}
}

func commandPlaceMarker(g *Game) {
	g.serverCommandsCh <- server.NewCMPlaceMarker(g.cursorY, g.cursorX)
}

func commandQuit(g *Game) {
	g.quit = true
}
