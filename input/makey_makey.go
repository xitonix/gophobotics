package input

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// MakeyMakey implements the Source interface and provides keypress commands from a MakeyMakey board
type MakeyMakey struct {
	commands  chan Command
	started   bool
	verbosity Verbosity
}

// NewMakeyMakey creates a new MakeyMakey source
func NewKMakeyMakey(verbosity Verbosity) *MakeyMakey {
	return &MakeyMakey{
		commands:  make(chan Command),
		verbosity: verbosity,
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
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)
	printHelp(false)
	for {
		ev := termbox.PollEvent()
		cmd := t.parseKey(ev.Key)
		if t.verbosity == VeryVerbose {
			fmt.Printf("KEY: %v, CH: %v, MODIFIER: %v, EVENT: %v\n", ev.Key, ev.Ch, ev.Mod, ev.Type)
		}
		t.commands <- cmd

		if cmd == Exit {
			t.started = false
			close(t.commands)
			return nil
		}
	}
}

func (t *MakeyMakey) parseKey(key termbox.Key) Command {
	cmd := None
	switch key {
	case termbox.MouseLeft, termbox.KeyCtrlC:
		cmd = Exit
	case termbox.KeyArrowUp:
		cmd = Forward
	case termbox.KeyArrowDown:
		cmd = Backward
	case termbox.KeyArrowLeft:
		cmd = Left
	case termbox.KeyArrowRight:
		cmd = Right

	case termbox.KeySpace:
		if !t.started {
			t.started = true
			cmd = TakeOff
		} else {
			t.started = false
			cmd = Land
		}
	}
	if t.verbosity >= Verbose {
		fmt.Printf("Command Triggered: %s\n", cmd)
	}

	return cmd
}
