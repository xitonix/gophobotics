package main

import (
	"log"
	"os/exec"

	"github.com/spf13/pflag"
	"github.com/xitonix/gophobotics/input"
	"github.com/xitonix/gophobotics/robot"
)

func main() {
	v := pflag.CountP("verbose", "v", "Enables verbose mode. You can enable extra verbosity by using -vv")
	maxMoves := pflag.IntP("max-moves", "m", 6, "Maximum number of allowed forward/backward/left/right moves")
	pflag.Parse()

	verbosity := input.ParseVerbosity(*v)
	source := input.NewKeyboard(verbosity)
	robo := robot.NewTello(30, *maxMoves, verbosity)

	mplayer := exec.Command("mplayer", "-fps", "60", "-")

	mplayerIn, err := mplayer.StdinPipe()
	if nil != err {
		panic(err)
	}

	if err := mplayer.Start(); err != nil {
		log.Fatal(err)
	}
	if err := robo.Video(mplayerIn); nil != err {
		log.Fatal(err)
	}

	go func() {
		err := robo.Connect(source)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		err = source.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		robo.MonitorTermination()
		_ = mplayer.Process.Kill()
	}()

	if err := mplayer.Wait(); err != nil {
		log.Printf("mplayer:%s\n", err)
	}
}
