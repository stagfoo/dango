package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
  "os/exec"
  kdl "github.com/sblinch/kdl-go"


	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

)

const inventoryFile = ".inventory"
const listHeight = 14

type model struct {
    choices  []string           // items on the to-do list
    cursor   int                // which to-do list item our cursor is pointing at
    selected map[int]struct{}   // which to-do items are selected
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <command> [arguments]")
		fmt.Println("Commands: pickup, putdown, inventory, lootdrop")
		return
	}

	switch os.Args[1] {
	case "pickup":
		pickup()
	case "putdown":
		putdown()
	case "inventory":
		inventory()
	case "lootdrop":
		lootdrop()
	default:
		fmt.Println("Invalid command")
	}
}

func getCurrentFolder() (string, error) {
  cmd := exec.Command("pwd")
	pwdBytes, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running pwd command: %v", err)
	}

	// Convert the byte slice to a string and trim any newline characters
	pwd := strings.TrimSpace(string(pwdBytes))

  return pwd, nil
}

func listFilesInFolder() ([]string, error) {
	pwd, err := getCurrentFolder()

	// Read the directory contents
	files, err := os.ReadDir(pwd)
  if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}
  
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func pickup() {

	files, err := listFilesInFolder()
	if err != nil {
		fmt.Println("Error running ls:", err)
		return
	}
  // get inventory files from the file
	program := tea.NewProgram(initialModel(files))
	if err, _ := program.Run(); err != nil {
		os.Exit(1)
	}
}

func putdown() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		filePath := scanner.Text()
		if filePath != "" {
			input, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("File not found:", filePath)
				continue
			}

			fileName := filepath.Base(filePath)
			err = os.WriteFile(fileName, input, 0644)
			if err != nil {
				fmt.Println("Error copying file:", err)
				continue
			}

			fmt.Println("You put down:", filePath)
			removeFromInventory(filePath)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}

func inventory() {
	content, err := os.ReadFile(inventoryFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Inventory is empty.")
			return
		}
		fmt.Println("Error reading inventory file:", err)
		return
	}

	files := strings.Split(string(content), "\n")
	p := tea.NewProgram(initialModel(files))
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %s\n", err)
		os.Exit(1)
	}
}

func lootdrop() {
	err := os.Remove(inventoryFile)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error clearing inventory file:", err)
		return
	}
	fmt.Println("All files have been removed from the inventory.")
}

func addToInventory(filePath string) {
  pwd, err := getCurrentFolder()
  input, err := os.ReadFile(inventoryFile)
	if err != nil {
		fmt.Println("Error reading inventory file:", err)
		return
	}

	lines := strings.Split(string(input), "\n")
	updated := false
	for i, line := range lines {
		if strings.TrimSpace(line) == strings.TrimSpace(filePath) {
			lines[i] = filePath // Update the line if it already exists
			updated = true
			break
		}
	}

	if !updated {
		lines = append(lines, pwd + "/" + filePath) // Add new line if it doesn't exist
	}

	output := strings.Join(lines, "\n")

	err = os.WriteFile(inventoryFile, []byte(output), 0644)
	if err != nil {
		fmt.Println("Error writing inventory file:", err)
	}

}

func removeFromInventory(filePath string) {
	input, err := os.ReadFile(inventoryFile)
	if err != nil {
		fmt.Println("Error reading inventory file:", err)
		return
	}

	lines := strings.Split(string(input), "\n")
	var output []string
	for _, line := range lines {
		if strings.TrimSpace(line) != strings.TrimSpace(filePath) {
			output = append(output, line)
		}
	}

	err = os.WriteFile(inventoryFile, []byte(strings.Join(output, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing inventory file:", err)
	}
}

func initialModel(choices []string) model {
	items := []list.Item{}
	for _, choice := range choices {
		if choice != ""  {
			items = append(items, listItem(strings.TrimSpace(choice)))
		}
	}
	m := model{
		choices: choices,
    selected: make(map[int]struct{}),

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

        // These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit

        case "ctrl+e", "e":
            fmt.Println("inventory saved")
            for line := range m.selected {
              addToInventory(m.choices[line])
            }

        // The "up" and "k" keys move the cursor up
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }

        // The "down" and "j" keys move the cursor down
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }

       //  The "enter" key and the spacebar (a literal space) toggle
       //  the selected state for the item that the cursor is pointing at.
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        }
      }



	return m, nil
}

func (m model) View() string {
    // The header
    s := "What file do you want to pickup?\n\n"

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
            checked = "🎒" // selected!
        }

        // Render the row
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }

    // The footer
    s += "\nPress q to quit.\n"

    // Send the UI for rendering
    return s
}
