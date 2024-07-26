package main

import (
	"os"
  "fmt"
	"inventory-cli/internal"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := internal.NewModel()

	p := tea.NewProgram(m)

	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}

