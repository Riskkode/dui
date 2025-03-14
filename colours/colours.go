package colours

import (
	"fmt"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

// Print prints text in a specific color.
func Print(color, text string) {
	Fprint(color, text)
}

// Fprint writes colored text.
func Fprint(color, text string) {
	fmt.Println(color, text, Reset)
}

func SFprint(color, text string) string {
	return fmt.Sprintln(color, text, Reset)
}

// FprintParts allows mixing colors within a single output.
func FprintParts(parts, colors []string) {
	for i := range parts {
		Fprint(colors[i], parts[i])
	}
}

func Err(text string) {
	Fprint(Red, text)
}

func Info(text string) {
	Fprint(Blue, text)
}

func Ok(text string) {
	Fprint(Green, text)
}

func SErr(text string) string {
	return SFprint(Red, text)
}

func SInfo(text string) string {
	return SFprint(Blue, text)
}

func SOk(text string) string {
	return SFprint(Green, text)
}
