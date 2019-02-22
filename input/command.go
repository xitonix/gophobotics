package input

type Command int8

const (
	Start Command = iota
	Stop

	Up
	Down
	Forward
	Backward
	Left
	Right

	FrontFlip
	BackFlip
	LeftFlip
	RightFlip
	Bounce
)

func (c Command) String() string {
	switch c {
	case Start:
		return "Start"
	case Stop:
		return "Stop"
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	case Forward:
		return "Forward"
	case Backward:
		return "Backward"
	case FrontFlip:
		return "FrontFlip"
	case BackFlip:
		return "BackFlip"
	case LeftFlip:
		return "LeftFlip"
	case RightFlip:
		return "RightFlip"
	case Bounce:
		return "Bounce"
	default:
		return "Unknown"
	}
}
