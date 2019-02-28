package input

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

// MakeyMakey implements the Source interface and provides keypress commands from a MakeyMakey board
type MakeyMakey struct {
	commands chan Command
	started  bool
	verbose  bool
}

// NewMakeyMakey creates a new MakeyMakey source
func NewKMakeyMakey(verbose bool) *MakeyMakey {
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
	printHelp(false)
	for {
		ev := termbox.PollEvent()
		exit := t.parseKey(ev.Key)
		if t.verbose {
			fmt.Printf("KEY: %v, CH: %v, MODIFIER: %v, EVENT: %v\n", ev.Key, ev.Ch, ev.Mod, ev.Type)
		}

		if exit {
			return nil
		}
	}
}

func (t *MakeyMakey) parseKey(key termbox.Key) bool {
	switch key {

	case termbox.MouseLeft, termbox.KeyCtrlC:
		t.started = false
		close(t.commands)
		return true

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
			t.commands <- TakeOff
		} else {
			t.started = false
			t.commands <- Land
		}
	}
	return false
}
