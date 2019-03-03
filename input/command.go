package input

type Command int8

const (
	None Command = iota
	TakeOff
	Land

	Up
	Down
	Forward
	Backward
	Left
	Right
	RotateRight
	RotateLeft

	FrontFlip
	BackFlip
	LeftFlip
	RightFlip
	Bounce

	Exit
)

func (c Command) String() string {
	switch c {
	case TakeOff:
		return "Takeoff"
	case Land:
		return "Land"
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	case RotateLeft:
		return "RotateLeft"
	case RotateRight:
		return "RotateRight"
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
	case Exit:
		return "Exit"
	default:
		return "Unknown"
	}
}

func (c Command) IsRotation() bool {
	return c == RotateLeft || c == RotateRight
}

func (c Command) IsLandOrTakeoff() bool {
	return c == TakeOff || c == Land
}
