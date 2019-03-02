package input

type Verbosity int8

const (
	NonVerbose Verbosity = iota
	Verbose
	VeryVerbose
)

func ParseVerbosity(v int) Verbosity {
	if v >= 2 {
		return VeryVerbose
	}
	if v == 1 {
		return Verbose
	}
	return NonVerbose
}
