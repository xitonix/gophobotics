package output

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

// Video encapsulates video processing functionality
type Video struct {
	drone *tello.Driver
	wg    sync.WaitGroup

	close  chan interface{}
	closed bool
}

// NewVideo creates a new Video type
func NewVideo(drone *tello.Driver) *Video {
	return &Video{
		drone: drone,
		close: make(chan interface{}),
	}
}

func (v *Video) Stop() {
	_ = v.drone.On(tello.ConnectedEvent, func(data interface{}) {})
	_ = v.drone.On(tello.VideoFrameEvent, func(data interface{}) {})
	v.closed = true
	close(v.close)
	v.wg.Wait()
}

func (v *Video) Capture() error {
	mplayer, err := v.startVideo()
	if err != nil {
		return err
	}

	go func() {
		<-v.close
		_ = mplayer.Process.Kill()
	}()

	v.wg.Add(1)
	go func() {
		defer v.wg.Done()
		_ = mplayer.Wait()
	}()

	return nil
}

func (v *Video) startVideo() (*exec.Cmd, error) {
	mplayer := exec.Command("mplayer", "-fps", "60", "-")
	videoBuffer, err := mplayer.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := v.initialiseVideo(videoBuffer); err != nil {
		return nil, err
	}

	if err := mplayer.Start(); err != nil {
		return nil, err
	}

	return mplayer, nil
}

func (v *Video) initialiseVideo(output io.WriteCloser) error {
	err := v.drone.On(tello.ConnectedEvent, func(data interface{}) {
		err := v.drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
		if err != nil {
			fmt.Printf("Video: failed to set the encoder rate: %s\n", err)
			return
		}

		// it needs to send `StartVideo` to the drone every 100ms
		gobot.Every(100*time.Millisecond, func() {
			if !v.closed {
				if err := v.drone.StartVideo(); nil != err {
					fmt.Printf("failed to start video on drone: %s\n", err)
				}
			}
		})
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to video connection events: %s", err)
	}

	err = v.drone.On(tello.VideoFrameEvent, func(data interface{}) {
		pkt := data.([]byte)
		if !v.closed && len(pkt) > 0 {
			if _, err := output.Write(pkt); err != nil {
				fmt.Printf("Render: %s\n", err)
			}
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to video frame events: %s", err)
	}

	return nil
}
