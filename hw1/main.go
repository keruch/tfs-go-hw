package main

import "fmt"

type (
	Color     = int
	CharStyle = int
	Attribute int
	Settings  map[Attribute]int
	Mod       func(Settings) Settings
)

const (
	Size Attribute = iota
	Char
	Foreground
	Background
	Style
)

const (
	BgBlack Color = 40 + iota
	BgRed
	BgGreen
	BgBrown
	BgBlue
	BgPurple
	BgCyan
	BgGray
)

const (
	FgBlack Color = 30 + iota
	FgRed
	FgGreen
	FgBrown
	FgBlue
	FgPurple
	FgCyan
	FgGray
)

const (
	Normal CharStyle = iota
	Bold
	Underlined
	Blinking
	Reverse
)

func printSandGlass(mods ...Mod) {
	// default values of all attributes
	settings := Settings{
		Size:       8,
		Char:       '#',
		Foreground: FgRed,
		Background: BgGreen,
		Style:      Bold,
	}

	// apply functions to settings
	for _, mod := range mods {
		settings = mod(settings)
	}

	for lineNumber := 0; lineNumber < settings[Size]; lineNumber++ {
		for columnNumber := 0; columnNumber < settings[Size]; columnNumber++ {
			var symbol rune
			if (lineNumber == 0 || lineNumber == settings[Size]-1) ||
				(columnNumber == lineNumber || columnNumber == settings[Size]-lineNumber-1) {
				symbol = rune(settings[Char])
			} else {
				symbol = ' '
			}
			fmt.Printf("\033[%dm\033[%dm\033[%dm%c\033[0m\033[0m\033[0m",
				settings[Foreground], settings[Background], settings[Style], symbol)
		}
		fmt.Println()
	}
}

func main() {
	printSandGlass(setAttribute(Size, 5), setAttribute(Background, BgBlue),
		setAttribute(Foreground, FgRed), setAttribute(Char, 'F'), setAttribute(Style, Bold))
}

func setAttribute(attr Attribute, value int) Mod {
	return func(settings Settings) Settings {
		settings[attr] = value
		return settings
	}
}
