package input

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// MakeyMakey implements the Source interface and provides commands from a Makey Makey board
type MakeyMakey struct {
	commands chan Command
	started  bool
	verbose  bool
}

// NewMakeyMakey creates a new MakeyMakey source
func NewMakeyMakey(verbose bool) *MakeyMakey {
	return &MakeyMakey{
		commands: make(chan Command),
		verbose:  verbose,
	}
}

func (t *MakeyMakey) Commands() <-chan Command {
	return t.commands
}

func (t *MakeyMakey) Start() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)
	for {
		ev := termbox.PollEvent()
		switch ev.Key {
		case termbox.MouseLeft:
			t.started = false
			close(t.commands)
			return nil
		case termbox.KeyArrowUp:
			t.commands <- Forward
		case termbox.KeyArrowDown:
			t.commands <- Backward
		case termbox.KeyArrowLeft:
			t.commands <- Left
		case termbox.KeyArrowRight:
			t.commands <- Right
		case termbox.KeySpace:
			if !t.started {
				t.started = true
				t.commands <- Start
			} else {
				t.started = false
				t.commands <- Stop
			}
		}

		if t.verbose {
			fmt.Printf("KEY: %v, CH: %v, MODIFIER: %v, EVENT: %v\n", ev.Key, ev.Ch, ev.Mod, ev.Type)
		}
	}
}
