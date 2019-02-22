package main

import (
	"log"
	"os/exec"
	"sync"

	"git.campmon.com/golang/corekit/proc"
	"github.com/spf13/pflag"
	"github.com/xitonix/gophobotics/input"
	"github.com/xitonix/gophobotics/robot"
)

func main() {
	verbose := pflag.BoolP("verbose", "v", false, "Enables verbose mode")
	pflag.Parse()
	source := input.NewKeyboard(*verbose)
	robo := robot.NewTello(50)
	mplayer := exec.Command("mplayer", "-fps", "60", "-")

	mplayerIn, err := mplayer.StdinPipe()
	if nil != err {
		panic(err)
	}

	if err := mplayer.Start(); err != nil {
		panic(err)
	}
	if err := robo.Video(mplayerIn); nil != err {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	if err := mplayer.Wait(); nil != err {
		log.Printf("mplayer error:%s\n", err)
	}
	proc.WaitForTermination()
}
