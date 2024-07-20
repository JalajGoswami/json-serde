package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var LogInfo = log.New(os.Stdout, getColored("blue", "Info:"), 0).Println
var LogSuccess = log.New(os.Stdout, getColored("green", "Success:"), 0).Println
var LogError = log.New(os.Stderr, getColored("red", "Error:"), 0).Println
var LogWarning = log.New(os.Stdout, getColored("yellow", "Warn:"), 0).Println

func getColored(color string, v any) string {
	var colorMap = map[string]string{
		"red":     "\x1b[31m\x1b[47m%v\x1b[0m ",
		"green":   "\x1b[32m\x1b[47m%v\x1b[0m ",
		"yellow":  "\x1b[33m\x1b[47m%v\x1b[0m ",
		"blue":    "\x1b[34m\x1b[47m%v\x1b[0m ",
		"magenta": "\x1b[35m\x1b[47m%v\x1b[0m ",
	}
	colored, ok := colorMap[color]
	if ok {
		return fmt.Sprintf(colored, v)
	}
	return fmt.Sprintf("%v", v)
}

func PrintLogo() {
	fmt.Println()
	fmt.Println(getColored("blue", strings.Repeat("-", 30)))
	fmt.Println(getColored("magenta", strings.Repeat(" ", 10)+"JSON Serde"+strings.Repeat(" ", 10)))
	fmt.Println(getColored("blue", strings.Repeat("-", 30)))
	fmt.Print("\n\n")
}
func OpenFile(paths ...*string) *os.File {
	var filePath string
	for _, path := range paths {
		if *path != "" {
			filePath = *path
		}
	}
	if filePath == "" {
		LogError("file path is missing! see --help")
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		LogError("error in opening file,", err)
		os.Exit(2)
	}

	stat, err := file.Stat()
	if err != nil {
		LogError("error in reading file,", err)
		os.Exit(3)
	}

	if stat.IsDir() {
		LogError("path provided must be of file, got directory")
		os.Exit(4)
	}
	return file
}
