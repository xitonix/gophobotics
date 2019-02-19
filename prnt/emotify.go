package prnt

import "fmt"

// Emotifier adds emoticons to the message if enabled
type Emotifier struct {
	enabled bool
}

// NewEmotifier creates a new emotifier
func NewEmotifier(enabled bool) Emotifier {
	return Emotifier{enabled: enabled}
}

//Println writes a line prefixed with `emoticon ` into stdout
func (e Emotifier) Println(emoticon, message string) {
	if e.enabled {
		fmt.Printf("%s %s\n", emoticon, message)
	} else {
		fmt.Println(message)
	}
}

//Printf writes a formatted string prefixed with `emoticon ` into stdout
func (e Emotifier) Printf(emoticon, format string, args ...interface{}) {
	if e.enabled {
		msg := fmt.Sprintf(format, args...)
		fmt.Printf("%s %s", emoticon, msg)
	} else {
		fmt.Printf(format, args...)
	}
}
