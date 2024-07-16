package main

import (
	"flag"
	"fmt"
	"loro-tui/core"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func main() {
	serverEndpoint := flag.String("server", "", "Chat server endpoint (required)")
	flag.Parse()

	if *serverEndpoint == "" {
		fmt.Println("Error: The -server flag is required")
		flag.Usage()
		os.Exit(1)
	}

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	if _, err := tea.NewProgram(core.NewModel(width, height, *serverEndpoint)).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
