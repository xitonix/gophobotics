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

// Video setup video feeds
// it need to be called before you connect to other source
func (t *Tello) Video(output io.WriteCloser) error {
	if nil == output {
		return nil
	}
	t.drone.On(tello.ConnectedEvent, func(data interface{}) {
		t.drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
		t.drone.StartVideo()
		// it need to send `StartVideo` to the drone every 100ms
		gobot.Every(100*time.Millisecond, func() {
			if err := t.drone.StartVideo(); nil != err {
				fmt.Printf("fail to start video on drone:%s\n", err)
			}
		})
	})

	t.drone.On(tello.VideoFrameEvent, func(data interface{}) {
		pkt := data.([]byte)
		if len(pkt) > 0 {
			if _, err := output.Write(pkt); err != nil {
				fmt.Printf("err:%s\n", err)
			}
		}
	})
	return nil
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
