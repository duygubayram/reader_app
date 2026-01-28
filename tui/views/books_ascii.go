package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func bookStyle(color lipgloss.Color, selected bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Foreground(color).
		Border(lipgloss.ThickBorder()).
		BorderForeground(color).
		Width(5).
		Height(7).
		Align(lipgloss.Center)

	if selected {
		style = style.
			BorderForeground(lipgloss.Color("#FFFFFF")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3A86FF")).
			Bold(true)
	}

	return style
}

func renderBook(name string, color lipgloss.Color, selected bool) string {
	name = strings.ToUpper(name)
	runes := []rune(name)

	lines := make([]string, 0, 5)
	for i := 0; i < len(runes) && i < 5; i++ {
		lines = append(lines, string(runes[i]))
	}
	for len(lines) < 5 {
		lines = append(lines, " ")
	}

	content := strings.Join(lines, "\n")

	return bookStyle(color, selected).Render(content)
}

