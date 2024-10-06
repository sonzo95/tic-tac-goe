package client

import "testing"

type cursorTest struct {
	cX, cY, wantX, wantY int
}

func TestCommandUp(t *testing.T) {
	tests := []cursorTest{
		{1, 1, 1, 0},
		{1, 0, 1, 0},
		{1, 2, 1, 1},
	}

	t.Run("moves cursor up", func(t *testing.T) {
		for _, test := range tests {
			g := Game{}
			g.cursorX = test.cX
			g.cursorY = test.cY

			commandUp(&g)
			assertCursor(t, &g, test.wantX, test.wantY)
		}
	})
}

func TestCommandDown(t *testing.T) {
	tests := []cursorTest{
		{1, 1, 1, 2},
		{1, 0, 1, 1},
		{1, 2, 1, 2},
	}

	t.Run("moves cursor down", func(t *testing.T) {
		for _, test := range tests {
			g := Game{}
			g.cursorX = test.cX
			g.cursorY = test.cY

			commandDown(&g)
			assertCursor(t, &g, test.wantX, test.wantY)
		}
	})
}

func TestCommandLeft(t *testing.T) {
	tests := []cursorTest{
		{2, 1, 1, 1},
		{1, 1, 0, 1},
		{0, 1, 0, 1},
	}

	t.Run("moves cursor left", func(t *testing.T) {
		for _, test := range tests {
			g := Game{}
			g.cursorX = test.cX
			g.cursorY = test.cY

			commandLeft(&g)
			assertCursor(t, &g, test.wantX, test.wantY)
		}
	})
}

func TestCommandRight(t *testing.T) {
	tests := []cursorTest{
		{0, 1, 1, 1},
		{1, 1, 2, 1},
		{2, 1, 2, 1},
	}

	t.Run("moves cursor right", func(t *testing.T) {
		for _, test := range tests {
			g := Game{}
			g.cursorX = test.cX
			g.cursorY = test.cY

			commandRight(&g)
			assertCursor(t, &g, test.wantX, test.wantY)
		}
	})
}

func TestCommandPlaceMarker(t *testing.T) {
	// TODO: test some channel or thing that writes onto websocket
	t.Run("TBD", func(t *testing.T) {
		g := Game{}
		commandPlaceMarker(&g)
	})
}

func TestCommandQuit(t *testing.T) {
	t.Run("quits game", func(t *testing.T) {
		g := Game{}
		commandQuit(&g)
		if !g.quit {
			t.Errorf("expected quit to be true")
		}
	})
}

func assertCursor(t testing.TB, got *Game, wantX, wantY int) {
	t.Helper()
	if got.cursorX != wantX || got.cursorY != wantY {
		t.Errorf("expected cursor to be at %d:%d, found it at %d:%d", wantX, wantY, got.cursorX, got.cursorY)
	}
}
