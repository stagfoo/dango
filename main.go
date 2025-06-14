package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	internal "dango/internal"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pelletier/go-toml"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	quitting string
}

var ViewName string
var Message string
var PreviousSelection string = ""

func main() {
	ensureDBExists()

	if len(os.Args) < 2 {
		fmt.Println("Usage: dango <command>")
		return
	}

	command := os.Args[1]

	switch command {
	case "stick":
		ViewName = "pickup"
		files := readInput()
		MyProgram := tea.NewProgram(initialModel(files))
		internal.ViewDBItems(MyProgram)
	case "stuck":
		ViewName = "list"
		db := internal.ViewDB(internal.Path)
		MyProgram := tea.NewProgram(initialModel(db.Items))
		internal.ViewDBItems(MyProgram)
	case "drop":
		internal.RemoveAllItems()
	case "serve":
		internal.OutputSelectedItems()
	case "box":
		// prompt user to create new toml
		// list dangos to edit
		//
	default:
		//TODO add option for multiple dangos and make this i select of dangos
		// alt command will be box
		fmt.Println("🍡 No command provided <stick |stuck | serve | drop>")
	}
}

func ensureDBExists() {
	if _, err := os.Stat(internal.Path); os.IsNotExist(err) {
		dir := filepath.Dir(internal.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
			os.Exit(1)
		}

		emptyDB := internal.MyDB{Items: []string{}}
		data, err := toml.Marshal(emptyDB)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling empty DB: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(internal.Path, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating dango.toml: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Created new dango.toml file.")
	}
}

func readInput() []string {
	var input []string
	HasPipedInput := internal.HasInputFromPipe()
	if !HasPipedInput {
		// If no input from stdin, fallback to listing files in the folder
		files, err := internal.ListDirectoryContents(false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting files in folder: %v\n", err)
			os.Exit(1)
		}
		input = files
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input = append(input, scanner.Text())
		}
	}
	return input
}

func initialModel(choices []string) model {
	selected := make(map[int]struct{})
	for cursor, choice := range choices {
		if choice != "" {
			if internal.FindInDB(choice) != "" {
				selected[cursor] = struct{}{}
			}
		}
	}
	return model{
		choices:  choices,
		selected: selected,
	}
}

func (m model) Init() tea.Cmd {
	ViewName = "pickup"
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// This is keypress message there are other message types in bubbletea
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			_, ok := m.selected[m.cursor]
			if ok {
				tea.Println(m.selected[m.cursor])
				internal.CopyToClipboard(m.choices[m.cursor])
				return m, tea.Quit
			}
		// These keys should exit the program.
		case "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			Message = ""

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
			Message = ""

			//  The "enter" key and the spacebar (a literal space) toggle
			//  the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				internal.RemoveFromDB(m.choices[m.cursor])
				delete(m.selected, m.cursor)
				Message = "Removed"
			} else {
				internal.AddToDB(m.choices[m.cursor])
				m.selected[m.cursor] = struct{}{}
				Message = "Picked up"
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	// The header
	s := "[space] to add/remove.\n"
	s += "c to copy.\n"
	s += "q to close.\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := "✨" // not selected
		if _, ok := m.selected[i]; ok {
			checked = "🍡" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += "\n"
	s += Message
	// Send the UI for rendering
	return s
}
