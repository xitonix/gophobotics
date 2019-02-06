package input

// Source is an interface to provide input commands to the robot
type Source interface {
	Commands() <-chan Command
}
