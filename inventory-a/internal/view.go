package internal

import (
	"github.com/charmbracelet/lipgloss"
	//"github.com/charmbracelet/bubbletea"
)

func (m Model) View() string {
	var s string
	for i, item := range m.Items {
		cursor := " " // no cursor
		if m.Selected == i {
			cursor = ">" // cursor
		}
		s += lipgloss.NewStyle().Bold(true).Render(cursor + " " + item.Name) + "\n"
	}
	return s
}

