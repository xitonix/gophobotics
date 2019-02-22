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
	source := input.NewKeyboard(*verbose)

	robot := robot.NewTello(50, true)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range robot.Errors() {
			fmt.Printf("Err: %s\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := robot.Connect(source)
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
