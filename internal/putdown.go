package internal

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

//NOTE Exports have to have capitalized names

// this function is weird, why is it needed?
func Putdown(program *tea.Program) []string {

	_, runError := program.Run()
	if runError != nil {
		fmt.Printf("Error running program: %s\n", runError)
	}

	return ViewDB(Path).Items
}
