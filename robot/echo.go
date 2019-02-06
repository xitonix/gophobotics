package robot

import (
	"fmt"

	"go.xitonix.io/gophobotics/input"
)

// Echo is a robot that simply prints every input command to standard output
type Echo struct {
	errors chan error
}

func NewEcho() *Echo {
	return &Echo{
		errors: make(chan error),
	}
}

func (e *Echo) Errors() <-chan error {
	return e.errors
}

// Connect is a blocking call which blocks until the context has been cancelled
func (e *Echo) Connect(source input.Source) error {
	defer close(e.errors)
	for cmd := range source.Commands() {
		fmt.Printf("Command Received: %s\n", cmd)
	}
	return nil
}
