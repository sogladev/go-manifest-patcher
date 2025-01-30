package util

import "runtime"

// ANSI color codes
const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Cyan   = "\033[36m"
)

// Colorize wraps the given text with the specified color code
func Colorize(text, color string) string {
	if runtime.GOOS == "windows" {
		return text // Return plain text on Windows
	}
	return color + text + Reset
}

// Color functions for convenience
func ColorGreen(text string) string {
	if runtime.GOOS == "windows" {
		return text // Return plain text on Windows
	}
	return Colorize(text, Green)
}

func ColorYellow(text string) string {
	if runtime.GOOS == "windows" {
		return text // Return plain text on Windows
	}
	return Colorize(text, Yellow)
}

func ColorRed(text string) string {
	if runtime.GOOS == "windows" {
		return text // Return plain text on Windows
	}
	return Colorize(text, Red)
}

func ColorCyan(text string) string {
	if runtime.GOOS == "windows" {
		return text // Return plain text on Windows
	}
	return Colorize(text, Cyan)
}
