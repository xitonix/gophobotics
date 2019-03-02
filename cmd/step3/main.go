package main

import (
	"log"
	"sync"

	"github.com/spf13/pflag"
	"github.com/xitonix/gophobotics/input"
	"github.com/xitonix/gophobotics/robot"
)

func main() {
	v := pflag.CountP("verbose", "v", "Enables verbose mode. You can enable extra verbosity by using -vv")
	pflag.Parse()

	source := input.NewKeyboard(input.ParseVerbosity(*v))
	robo := robot.NewEcho()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := robo.Connect(source)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := source.Start()
	if err != nil {
		log.Fatal(err)
	}
}
