package main

import "fmt"

type Color int
type CharStyle int
type SandGlassMod func(size int, char rune, str string) (int, rune, string)

const (
	Black Color = iota
	Red
	Green
	Brown
	Blue
	Purple
	Cyan
	Gray
)

const (
	Normal CharStyle = iota
	Bold
	Underlined
	Blinking
	Reverse
)

func printSandGlass(mods ...SandGlassMod) {
	// default values for size and char
	size, char := 8, '#'

	// apply functions to size and char
	for _, mod := range mods {
		size, char, _ = mod(size, char, "")
	}

	for lineNumber := 0; lineNumber < size; lineNumber++ {
		for columnNumber := 0; columnNumber < size; columnNumber++ {
			var symbol rune
			if (lineNumber == 0 || lineNumber == size-1) ||
				(columnNumber == lineNumber || columnNumber == size-lineNumber-1) {
				symbol = char
			} else {
				symbol = ' '
			}
			formatString := string(symbol)

			// apply functions to format string
			for _, mod := range mods {
				_, _, formatString = mod(0, 0, formatString)
			}
			fmt.Print(formatString)
		}
		fmt.Println()
	}
}

func main() {
	printSandGlass(setSandGlassChar('H'), setSandGlassSize(5),
		setCharStyle(Reverse), setBackgroundColor(Red), setForegroundColor(Blue))
}

func setForegroundColor(color Color) SandGlassMod {
	color += 30
	return func(size int, char rune, s string) (int, rune, string) {
		return size, char, fmt.Sprintf("\033[%dm%s\033[0m", color, s)
	}
}

func setBackgroundColor(color Color) SandGlassMod {
	color += 40
	return func(size int, char rune, s string) (int, rune, string) {
		return size, char, fmt.Sprintf("\033[%dm%s\033[0m", color, s)
	}
}

func setCharStyle(style CharStyle) SandGlassMod {
	return func(size int, char rune, s string) (int, rune, string) {
		return size, char, fmt.Sprintf("\033[%dm%s\033[0m", style, s)
	}
}

func setSandGlassSize(inputSize int) SandGlassMod {
	return func(size int, char rune, s string) (int, rune, string) {
		return inputSize, char, s
	}
}

func setSandGlassChar(inputChar rune) SandGlassMod {
	return func(size int, char rune, s string) (int, rune, string) {
		return size, inputChar, s
	}
}
