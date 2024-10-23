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
}

var ViewName string
var Message string

func main() {
	ensureDBExists()

	if len(os.Args) < 2 {
		fmt.Println("Usage: dango <command>")
		return
	}

	command := os.Args[1]

	switch command {
	case "pickup":
		ViewName = "pickup"
		files := readInput()
		MyProgram := tea.NewProgram(initialModel(files))
		internal.Pickup(MyProgram)
	case "list":
		ViewName = "list"
		db := internal.ViewDB(internal.Path)
		MyProgram := tea.NewProgram(initialModel(db.Items))
		internal.Putdown(MyProgram)
	case "drop":
		lootdrop()
	case "output":
		outputSelectedItems()
  case "show":
    listDirectoriesAndFiles()
	default:
		fmt.Println("🍡")
	}
}

func listDirectoriesAndFiles() error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Use filepath.Glob to get all files and directories in the current directory
	files, err := filepath.Glob(filepath.Join(currentDir, "*"))
	if err != nil {
		return err
	}

	// Print each file/directory path to stdout
	for _, file := range files {
		fmt.Println(file)
	}
	return nil
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
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	if len(input) == 0 {
		// If no input from stdin, fallback to listing files in the folder
		files, err := internal.ListFilesInFolder()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting files in folder: %v\n", err)
			os.Exit(1)
		}
		return files
	}
	return input
}

func outputSelectedItems() {
	db := internal.ViewDB(internal.Path)
	for _, item := range db.Items {
		fmt.Println(item)
	}
}

func lootdrop() {
	err := os.Remove(internal.Path)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error clearing dango file: %v\n", err)
		return
	}
	fmt.Println("All files have been removed from the dango.")
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

type listItem string

func (i listItem) FilterValue() string { return string(i) }

func (i listItem) Title() string { return string(i) }

func (i listItem) Description() string { return "" }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "c":
			_, ok := m.selected[m.cursor]
			if ok {
				internal.CopyToClipboard(m.choices[m.cursor])
				Message = "Copied to clipboard."
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
	s := "What file do you want to pickup?\n"
	s += "Press space to add.\n"
	s += "Press c to copy.\n"
	s += "Press q to close.\n\n"

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
