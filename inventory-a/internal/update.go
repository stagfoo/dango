package internal

import (
	bubbletea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case Msg:
		switch msg {
		case LoadItems:
			// Load items from KDL file
		case AddItem:
			// Add item to the model and KDL file
		case RemoveItem:
			// Remove item from the model and KDL file
		case ClearItems:
			// Clear all items from the model and KDL file
		case SelectItem:
			// Select item for piping
		}
	}
	return m, nil
}

