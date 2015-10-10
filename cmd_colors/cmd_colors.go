package cmd_colors

import (
	"github.com/mgutz/ansi"
	"fmt"
)

var greenFormat func(string) string = ansi.ColorFunc("green+b+h")
var redFormat func(string) string = ansi.ColorFunc("red+b+h")
var blueFormat func(string) string = ansi.ColorFunc("blue+b+h")
var resetFormat string = ansi.ColorCode("reset")

var Info = func(message string) {
	fmt.Println("\r" + greenFormat(message) + resetFormat)
}

var InfoInline = func(message string) {
	fmt.Print("\r" + blueFormat(message) + resetFormat)
}

var Warn = func(message string) {
	fmt.Println(redFormat(message) + resetFormat)
}