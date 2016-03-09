package cmdcolors

import (
	"github.com/mgutz/ansi"
	"fmt"
)

var greenFormat = ansi.ColorFunc("green+b+h")
var redFormat = ansi.ColorFunc("red+b+h")
var blueFormat = ansi.ColorFunc("blue+b+h")
var resetFormat = ansi.ColorCode("reset")

var Info = func(message string) {
	fmt.Println("\r" + greenFormat(message) + resetFormat)
}

var InfoInline = func(message string) {
	fmt.Print("\r" + blueFormat(message) + resetFormat)
}

var Warn = func(message string) {
	fmt.Println(redFormat(message) + resetFormat)
}