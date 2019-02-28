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
	verbose := pflag.BoolP("verbose", "v", false, "Enables verbose mode")
	pflag.Parse()
	source := input.NewKMakeyMakey(*verbose)

	tello := robot.NewTello(50)

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
