package robot

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/xitonix/gophobotics/input"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type Tello struct {
	drone                  *tello.Driver
	move, maxNumberOfMoves int
	errors                 chan error
	done                   chan interface{}
	terminated             chan interface{}
	closed                 bool
	moves                  struct {
		forward int
		back    int
		left    int
		right   int
		up      int
		down    int
	}
	verbosity        input.Verbosity
	internalCommands chan input.Command
}

// NewTello creates a new Tello drone robot
func NewTello(move, maxNumberOfMoves int, verbosity input.Verbosity) *Tello {
	return &Tello{
		drone:            tello.NewDriver("8888"),
		move:             move,
		maxNumberOfMoves: maxNumberOfMoves,
		errors:           make(chan error),
		done:             make(chan interface{}),
		terminated:       make(chan interface{}),
		verbosity:        verbosity,
		internalCommands: make(chan input.Command, 1000),
	}
}

// Errors returns any errors occurred during the execution of a command.
// MAKE SURE you always read from this channel before calling the Connect method to avoid deadlock
func (t *Tello) Errors() <-chan error {
	return t.errors
}

// Video setup video feeds
// it need to be called before you connect to other source
func (t *Tello) Video(output io.WriteCloser) error {
	if nil == output {
		return nil
	}
	_ = t.drone.On(tello.ConnectedEvent, func(data interface{}) {
		_ = t.drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
		_ = t.drone.StartVideo()
		// it need to send `StartVideo` to the drone every 100ms
		gobot.Every(100*time.Millisecond, func() {
			if t.closed {
				return
			}
			if err := t.drone.StartVideo(); nil != err {
				fmt.Printf("fail to start video on drone:%s\n", err)
			}
		})
	})

	_ = t.drone.On(tello.VideoFrameEvent, func(data interface{}) {
		if t.closed {
			return
		}
		pkt := data.([]byte)
		if len(pkt) > 0 {
			if _, err := output.Write(pkt); err != nil {
				fmt.Printf("err:%s\n", err)
			}
		}
	})
	return nil
}

func (t *Tello) MonitorTermination() {
	<-t.done
}

// Connect establishes a new connection to the drone and blocks until the source's Commands channel is closed.
func (t *Tello) Connect(source input.Source) error {
	_ = t.drone.On(tello.FlightDataEvent, t.flightData)

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{t.drone})

	go t.filter(source.Commands())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !t.closed {
			select {
			case <-t.terminated:
				t.drone.Hover()
				t.printCommand(input.Exit)
				t.closed = true
				close(t.internalCommands)
				close(t.done)
			case cmd, more := <-t.internalCommands:
				if !more {
					t.closed = true
				}
				err, ignored := t.executeCommand(cmd)
				if err != nil {
					t.errors <- err
					continue
				}

				if !ignored {
					t.printCommand(cmd)
				}

				if cmd.IsLandOrTakeoff() || ignored {
					continue
				}

				time.Sleep(500 * time.Millisecond)
				if cmd.IsRotation() {
					t.drone.CeaseRotation()
				} else {
					t.drone.Hover()
				}
			}
		}

		err := t.drone.Halt()
		if err != nil {
			t.errors <- err
		}
		close(t.errors)
		return
	}()

	go func() {
		_ = robot.Start()
	}()

	wg.Wait()
	return nil
}

func (t *Tello) printCommand(command input.Command) {
	if t.verbosity >= input.Verbose {
		fmt.Printf("Drone: %s Command Received\n", command)
	}
}

func (t *Tello) flightData(s interface{}) {
	if fd, ok := s.(*tello.FlightData); ok {
		if fd.BatteryLow {
			fmt.Printf("Battery is low %d%%\n", fd.BatteryPercentage)
			time.Sleep(5 * time.Second)
		}
	}
}

func (t *Tello) executeCommand(command input.Command) (error, bool) {
	switch command {
	case input.TakeOff:
		return t.drone.TakeOff(), false
	case input.Land:
		return t.drone.Land(), false

	case input.Left:
		if t.isOverLimit(command) {
			return nil, true
		}
		return t.drone.Left(t.move), false
	case input.Right:
		if t.isOverLimit(command) {
			return nil, true
		}
		return t.drone.Right(t.move), false
	case input.Forward:
		if t.isOverLimit(command) {
			return nil, true
		}
		return t.drone.Forward(t.move), false
	case input.Backward:
		if t.isOverLimit(command) {
			return nil, true
		}
		return t.drone.Backward(t.move), false
	case input.RotateRight:
		return t.drone.Clockwise(t.move), false
	case input.RotateLeft:
		return t.drone.CounterClockwise(t.move), false

	case input.Up:
		if t.isOverLimit(command) {
			return nil, true
		}
		return t.drone.Up(t.move), false
	case input.Down:
		if t.isOverLimit(command) {
			return nil, true
		}
		return t.drone.Down(t.move), false

	case input.FrontFlip:
		return t.drone.FrontFlip(), false
	case input.BackFlip:
		return t.drone.BackFlip(), false
	case input.LeftFlip:
		return t.drone.LeftFlip(), false
	case input.RightFlip:
		return t.drone.RightFlip(), false

	case input.Bounce:
		return t.drone.Bounce(), false
	default:
		return nil, true
	}
}

func (t *Tello) isOverLimit(cmd input.Command) bool {
	if t.maxNumberOfMoves <= 0 {
		return false
	}
	var current int
	switch cmd {
	case input.Left:
		if t.moves.left < t.maxNumberOfMoves {
			t.moves.left++
		}
		if t.moves.right > 0 {
			t.moves.right--
		}
		current = t.moves.left

	case input.Right:
		if t.moves.right < t.maxNumberOfMoves {
			t.moves.right++
		}
		if t.moves.left > 0 {
			t.moves.left--
		}
		current = t.moves.right
	case input.Forward:
		if t.moves.forward < t.maxNumberOfMoves {
			t.moves.forward++
		}
		if t.moves.back > 0 {
			t.moves.back--
		}
		current = t.moves.forward
	case input.Backward:
		if t.moves.back < t.maxNumberOfMoves {
			t.moves.back++
		}
		if t.moves.forward > 0 {
			t.moves.forward--
		}
		current = t.moves.back
	case input.Up:
		if t.moves.up < t.maxNumberOfMoves {
			t.moves.up++
		}
		if t.moves.down > 0 {
			t.moves.down--
		}
		current = t.moves.up
	case input.Down:
		if t.moves.down < t.maxNumberOfMoves {
			t.moves.down++
		}
		if t.moves.up > 0 {
			t.moves.up--
		}
		current = t.moves.down
	default:
		return false
	}
	if t.verbosity >= input.VeryVerbose {
		fmt.Printf("Current Moves: %+v\n", t.moves)
	}
	return current >= t.maxNumberOfMoves
}

func (t *Tello) filter(commands <-chan input.Command) {
	for cmd := range commands {
		if cmd == input.Exit {
			close(t.terminated)
			return
		}
		t.internalCommands <- cmd
	}
}
