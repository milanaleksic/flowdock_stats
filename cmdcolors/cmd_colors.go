package cmdcolors

import (
	"fmt"

	"github.com/mgutz/ansi"
)

var greenFormat = ansi.ColorFunc("green+b+h")
var redFormat = ansi.ColorFunc("red+b+h")
var blueFormat = ansi.ColorFunc("blue+b+h")
var resetFormat = ansi.ColorCode("reset")

/*
Info shows a friendly green-colored message in the shell with an end-of-line character before it
*/
var Info = func(message string) {
	fmt.Println("\r" + greenFormat(message) + resetFormat)
}

/*
InfoInline shows a friendly blue-colored message in the shell with an end of line at the before it
*/
var InfoInline = func(message string) {
	fmt.Print("\r" + blueFormat(message) + resetFormat)
}

/*
Warn shows message in red text in the shell without an EOL
*/
var Warn = func(message string) {
	fmt.Println(redFormat(message) + resetFormat)
}
