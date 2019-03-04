package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/pflag"
	"github.com/xitonix/gophobotics/input"
	"github.com/xitonix/gophobotics/robot"
)

func main() {
	v := pflag.CountP("verbose", "v", "Enables verbose mode. You can enable extra verbosity by using -vv")
	maxMoves := pflag.IntP("max-moves", "m", 4, "Maximum number of allowed movements")
	pflag.Parse()
	verbosity := input.ParseVerbosity(*v)
	source := input.NewKeyboard(verbosity)

	tello := robot.NewTello(40, *maxMoves, verbosity)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range tello.Errors() {
			fmt.Printf("Err: %s\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := tello.Connect(source)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := source.Start()
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
