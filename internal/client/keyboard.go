package client

import "github.com/nsf/termbox-go"

var keymap = map[termbox.Key]Command{
	termbox.KeyArrowDown:  commandDown,
	termbox.KeyArrowUp:    commandUp,
	termbox.KeyArrowLeft:  commandLeft,
	termbox.KeyArrowRight: commandRight,
	termbox.KeyEnter:      commandPlaceMarker,
	termbox.KeyEsc:        commandQuit,
}

// Loops infinitely on keyboard events and writes them onto the channel
func ListenKeyboard(ch chan Command) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			cmd, found := keymap[ev.Key]
			if found {
				ch <- cmd
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
