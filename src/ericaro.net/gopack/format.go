package gopack

import (
	"fmt"
)
// use some escape character for vterm to pretty print text in a console
const (
	TERM_RESET     = 0
	TERM_BRIGHT    = 1
	TERM_DIM       = 2
	TERM_NULL      = 3
	TERM_UNDERLINE = 4
	TERM_BLINK     = 5
	TERM_REVERSE   = 7
	TERM_HIDDEN    = 8

	COLOR_BLACK   = 0
	COLOR_RED     = 1
	COLOR_GREEN   = 2
	COLOR_YELLOW  = 3
	COLOR_BLUE    = 4
	COLOR_MAGENTA = 5
	COLOR_CYAN    = 6
	COLOR_WHITE   = 7
	COLOR_DEFAULT = 9
)

//Defines some styles used in the command.
var (
	TitleStyle = PFormat{TERM_BRIGHT, COLOR_DEFAULT, COLOR_DEFAULT}
	ErrorStyle = PFormat{TERM_NULL, COLOR_RED, COLOR_DEFAULT}
	SuccessStyle = PFormat{TERM_NULL, COLOR_GREEN, COLOR_DEFAULT}
	NormalStyle = PFormat{TERM_NULL, COLOR_DEFAULT, COLOR_DEFAULT}
	
)

type PFormat struct {
	Attr, Foreground, Background int
}

func (f *PFormat) Printf(message string, v ...interface{}) {
	fmt.Print( f.Sprintf(message, v...) )
}

func (f *PFormat) Clear() {
	fmt.Printf("%c[1;1H%c[2J", 0x1B, 0x1B )
}

func (f *PFormat) Sprintf(message string, v ...interface{}) string{
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[%dm", 0x1B, f.Attr, f.Foreground+30, f.Background+40, fmt.Sprintf(message, v...), 0x1B, 0)
}

func (f *PFormat) PrintTriple(small, medium, large string) {
	f.Printf("       %-8s %-20s %-s\n", small, medium, large)
}
