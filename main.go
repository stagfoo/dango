package main

import (
	"fmt"
	"os"
	"strings"

	internal "dango/internal"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

var ViewName string
var Message string

func main() {
	//TODO add current files in folder to dango
	files, getFilesErr := internal.ListFilesInFolder()
	if getFilesErr != nil {
		fmt.Println("Error getting files in folder:", getFilesErr)
		return
	}

	switch os.Args[1] {
	case "pickup":
		ViewName = "pickup"
		MyProgram := tea.NewProgram(initialModel(files))
		internal.Pickup(MyProgram)
	case "list":
		ViewName = "list"
		db := internal.ViewDB(internal.Path)
		MyProgram := tea.NewProgram(initialModel(db.Items))
		internal.Putdown(MyProgram)
	case "lootdrop":
		lootdrop()
	default:
		fmt.Println("🍡")
		return
	}
}

func lootdrop() {
	err := os.Remove(internal.Path)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error clearing dango file:", err)
		return
	}
	fmt.Println("All files have been removed from the dango.")
}

func initialModel(choices []string) model {
	items := []list.Item{}
	selected := make(map[int]struct{})
	for cursor, choice := range choices {
		if choice != "" {
			if internal.FindInDB(choice) != "" {
				selected[cursor] = struct{}{}
				continue
			}
			_ = append(items, listItem(strings.TrimSpace(choice)))
		}
	}
	m := model{
		choices:  choices,
		selected: selected,
	}
	return m
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
