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
	commands  chan Command
	started   bool
	verbosity Verbosity
}

// NewKeyboard creates a new Keyboard source
func NewKeyboard(verbosity Verbosity) *Keyboard {
	return &Keyboard{
		commands:  make(chan Command),
		verbosity: verbosity,
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
	termbox.SetInputMode(termbox.InputAlt)
	printHelp(true)
	for {
		ev := termbox.PollEvent()
		if cmd := parseCharacter(ev.Ch); cmd != None {
			t.commands <- cmd
			continue
		}

		cmd := t.parseKey(ev.Key)
		if cmd == None {
			continue
		}
		t.commands <- cmd
		if t.verbosity == VeryVerbose {
			fmt.Printf("KEY: %v, CH: %v, MODIFIER: %v, EVENT: %v\n", ev.Key, ev.Ch, ev.Mod, ev.Type)
		}

		if cmd == Exit {
			t.started = false
			close(t.commands)
			_ = termbox.Clear(0, 0)
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

func printHelp(keyboard bool) {
	fmt.Printf("\nCONTROLS\n------------------------------\n")
	var mouse string
	if !keyboard {
		mouse = "/Left Click"
	}
	fmt.Printf("    CTRL + C%s: Emergency landing and EXIT\n\n", mouse)

	fmt.Printf("       SPACE: Takeoff/Land\n")
	fmt.Printf("    ARROW UP: Forward\n")
	fmt.Printf("  ARROW DOWN: Backward\n")
	fmt.Printf("  ARROW LEFT: Move left\n")
	fmt.Printf(" ARROW RIGHT: Move right\n")
	if keyboard {
		fmt.Printf("           L: Rotate Left\n")
		fmt.Printf("           R: Rotate Right\n")
		fmt.Printf("           U: UP\n")
		fmt.Printf("           D: Down\n")
		fmt.Printf("          F1: Front Flip (BE CAREFUL)\n")
		fmt.Printf("          F2: Back Flip (BE CAREFUL)\n")
		fmt.Printf("          F3: Right Flip (BE CAREFUL)\n")
		fmt.Printf("          F4: Left Flip (BE CAREFUL)\n")
		fmt.Printf("          F5: Bounce | Stop Bouncing (BE CAREFUL)\n")

	}
}

func (t *Keyboard) parseKey(key termbox.Key) Command {
	cmd := None
	switch key {

	case termbox.KeyCtrlC:
		cmd = Exit

	case termbox.KeyArrowUp:
		cmd = Forward
	case termbox.KeyArrowDown:
		cmd = Backward

	case termbox.KeyPgup:
		cmd = Up
	case termbox.KeyPgdn:
		cmd = Down
	case termbox.KeyArrowLeft:
		cmd = Left
	case termbox.KeyArrowRight:
		cmd = Right

	// Here goes the advanced move mappings

	case termbox.KeySpace:
		if !t.started {
			t.started = true
			cmd = TakeOff
		} else {
			t.started = false
			cmd = Land
		}
	}

	if t.verbosity >= Verbose && cmd != None {
		fmt.Printf("Command Triggered: %s\n", cmd)
	}
	return cmd
}
