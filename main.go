package main

import (
	"flag"
	"fmt"
	"loro-tui/core"
	"loro-tui/http_client"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func main() {
	url := flag.String("server", "", "Chat server url (required)")
	flag.Parse()

	if *url == "" {
		fmt.Println("Error: The -server flag is required")
		flag.Usage()
		os.Exit(1)
	}

	httpClient, err := http_client.NewClient(*url)
	if err != nil {
		fmt.Println("Server: " + err.Error())
		os.Exit(1)
	}

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	model := core.NewModel(width, height, httpClient)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
