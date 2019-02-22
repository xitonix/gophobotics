package robot

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/xitonix/gophobotics/input"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type Tello struct {
	drone       *tello.Driver
	move        int
	errors      chan error
	enableVideo bool
}

// NewTello creates a new Tello drone robot
func NewTello(move int, enableVideo bool) *Tello {
	return &Tello{
		drone:       tello.NewDriver("8888"),
		move:        move,
		errors:      make(chan error),
		enableVideo: enableVideo,
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
	if t.enableVideo {
		mplayer, err := t.startVideo()
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = mplayer.Wait()
		}()
	}

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
			if cmd.IsRotation() {
				t.drone.CeaseRotation()
			}
		}

		err := t.executeCommand(input.Land)
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
	case input.TakeOff:
		return t.drone.TakeOff()
	case input.Land:
		return t.drone.Land()
	case input.PalmLand:
		return t.drone.PalmLand()

	case input.Left:
		return t.drone.Left(t.move)
	case input.Right:
		return t.drone.Right(t.move)
	case input.RotateRight:
		return t.drone.Clockwise(t.move)
	case input.RotateLeft:
		return t.drone.CounterClockwise(t.move)
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

func (t *Tello) startVideo() (*exec.Cmd, error) {
	mplayer := exec.Command("mplayer", "-fps", "60", "-")
	videoBuffer, err := mplayer.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := t.initialiseVideo(videoBuffer); err != nil {
		return nil, err
	}

	if err := mplayer.Start(); err != nil {
		return nil, err
	}

	return mplayer, nil
}

func (t *Tello) initialiseVideo(output io.WriteCloser) error {
	err := t.drone.On(tello.ConnectedEvent, func(data interface{}) {
		err := t.drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
		if err != nil {
			t.errors <- fmt.Errorf("failed to set video encoder rate: %s", err)
		}
		err = t.drone.StartVideo()
		if err != nil {
			t.errors <- fmt.Errorf("failed start the video rate: %s", err)
		}
		// it needs to send `StartVideo` to the drone every 100ms
		gobot.Every(100*time.Millisecond, func() {
			if err := t.drone.StartVideo(); nil != err {
				fmt.Printf("failed to start video on drone:%s\n", err)
			}
		})
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to video connection events: %s", err)
	}

	err = t.drone.On(tello.VideoFrameEvent, func(data interface{}) {
		pkt := data.([]byte)
		if len(pkt) > 0 {
			if _, err := output.Write(pkt); err != nil {
				fmt.Printf("err:%s\n", err)
			}
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to video frame events: %s", err)
	}

	return nil
}
