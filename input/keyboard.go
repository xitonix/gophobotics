package input

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

var keyMap = map[rune]Command{
	117: Up,
	74:  Up,
	100: Down,
	68:  Down,
	108: RotateLeft,
	76:  RotateLeft,
	114: RotateRight,
	82:  RotateRight,
}

// Keyboard implements the Source interface and provides keypress commands from the keyboard
type Keyboard struct {
	commands chan Command
	started  bool
	verbose  bool
}

// NewKeyboard creates a new Keyboard source
func NewKeyboard(verbose bool) *Keyboard {
	return &Keyboard{
		commands: make(chan Command),
		verbose:  verbose,
	}
}

func (t *Keyboard) Commands() <-chan Command {
	return t.commands
}

func (t *Keyboard) Start() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt)
	for {
		ev := termbox.PollEvent()
		if cmd := parseCharacter(ev.Ch); cmd != None {
			t.commands <- cmd
			continue
		}

		exit := t.parseKey(ev.Key)

		if t.verbose {
			fmt.Printf("KEY: %v, CH: %v, MODIFIER: %v, EVENT: %v\n", ev.Key, ev.Ch, ev.Mod, ev.Type)
		}

		if exit {
			return nil
		}
	}
}

func parseCharacter(ch rune) Command {
	if cmd, ok := keyMap[ch]; ok {
		return cmd
	}
	return None
}

func (t *Keyboard) parseKey(key termbox.Key) bool {
	switch key {

	case termbox.KeyCtrlC:
		t.started = false
		close(t.commands)
		return true

	case termbox.KeyArrowUp:
		t.commands <- Forward
	case termbox.KeyArrowDown:
		t.commands <- Backward

	case termbox.KeyPgup:
		t.commands <- Up
	case termbox.KeyPgdn:
		t.commands <- Down
	case termbox.KeyArrowLeft:
		t.commands <- Left
	case termbox.KeyArrowRight:
		t.commands <- Right
	case termbox.KeyF1:
		t.commands <- FrontFlip
	case termbox.KeyF2:
		t.commands <- BackFlip
	case termbox.KeyF3:
		t.commands <- RightFlip
	case termbox.KeyF4:
		t.commands <- LeftFlip
	case termbox.KeyF5:
		t.commands <- Bounce

	case termbox.KeySpace, termbox.KeyCtrlSpace:
		if !t.started {
			t.started = true
			t.commands <- TakeOff
		} else {
			t.started = false
			if key == termbox.KeySpace {
				t.commands <- Land
			} else {
				t.commands <- PalmLand
			}
		}
	}
	return false
}
