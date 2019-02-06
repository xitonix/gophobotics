package robot

import (
	"sync"
	"time"

	"go.xitonix.io/gophobotics/input"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type Tello struct {
	drone  *tello.Driver
	move   int
	errors chan error
}

// NewTello creates a new Tello drone robot
func NewTello(move int) *Tello {
	return &Tello{
		drone:  tello.NewDriver("8888"),
		move:   move,
		errors: make(chan error),
	}
}

// Errors returns any errors occurred during the execution of a command.
// MAKE SURE you always read from this channel before calling the Connect method to avoid deadlock
func (t *Tello) Errors() <-chan error {
	return t.errors
}

// Connect establishes a new connection to the drone and blocks until the source's Commands channel is closed.
func (t *Tello) Connect(source input.Source) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cmd := range source.Commands() {
			err := t.executeCommand(cmd)
			if err != nil {
				t.errors <- err
			}
			time.Sleep(150 * time.Millisecond)
			t.drone.Hover()
		}

		err := t.executeCommand(input.Stop)
		if err != nil {
			t.errors <- err
		}
		close(t.errors)
		return
	}()

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{t.drone})

	err := robot.Start()
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func (t *Tello) executeCommand(command input.Command) error {
	switch command {
	case input.Start:
		return t.drone.TakeOff()
	case input.Stop:
		return t.drone.Land()

	case input.Left:
		return t.drone.Left(t.move)
	case input.Right:
		return t.drone.Right(t.move)
	case input.Forward:
		return t.drone.Forward(t.move)
	case input.Backward:
		return t.drone.Backward(t.move)
	case input.Up:
		return t.drone.Up(t.move)
	case input.Down:
		return t.drone.Down(t.move)

	case input.FrontFlip:
		return t.drone.FrontFlip()
	case input.BackFlip:
		return t.drone.BackFlip()
	case input.LeftFlip:
		return t.drone.LeftFlip()
	case input.RightFlip:
		return t.drone.RightFlip()

	case input.Bounce:
		return t.drone.Bounce()
	}
	return nil
}
