package main

import (
	"fmt"
	"os"
)

var tools map[string]func() = make(map[string]func())

func usage() {
	fmt.Printf("Usage: %v command [flags]\n\nValid commands:\n", os.Args[0])

	for tool := range tools {
		fmt.Printf("  %v\n", tool)
	}

	os.Exit(-1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	tool := os.Args[1]

	if f, ok := tools[tool]; !ok {
		usage()
	} else {
		args := []string{os.Args[0] + " " + tool}
		args = append(args, os.Args[2:]...)
		os.Args = args

		f()
	}
}
