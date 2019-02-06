package main

import (
	"log"
	"sync"

	"github.com/spf13/pflag"
	"go.xitonix.io/gophobotics/input"
	"go.xitonix.io/gophobotics/robot"
)

func main() {
	verbose := pflag.BoolP("verbose", "v", false, "Enables verbose mode")
	pflag.Parse()
	source := input.NewMakeyMakey(*verbose)
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
