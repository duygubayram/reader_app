package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var shelfPlank = lipgloss.NewStyle().
	Background(lipgloss.Color("#5C4033")).
	Foreground(lipgloss.Color("#F5DEB3")).
	Bold(true).
	Padding(0, 1)

var shelfFrame = lipgloss.NewStyle().
	Padding(1, 1)

func shelfColor(name string) lipgloss.Color {
	switch name {
	case "to_read":
		return lipgloss.Color("#5DA9E9")
	case "currently_reading":
		return lipgloss.Color("#F4D35E")
	case "read":
		return lipgloss.Color("#7AE582")
	default:
		return lipgloss.Color("#CCCCCC")
	}
}

func renderShelf(
	name string,
	books []string,
	shelfIndex int,
	selectedShelf int,
	selectedBook int,
) string {

	color := shelfColor(name)
	var renderedBooks []string

	for i, b := range books {
		isSelected := shelfIndex == selectedShelf && i == selectedBook
		renderedBooks = append(
			renderedBooks,
			renderBook(b, color, isSelected),
		)
	}

	bookRow := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		renderedBooks...,
	)

	title := shelfPlank.Render(" " + strings.ReplaceAll(name, "_", " ") + " ")

	return shelfFrame.Render(
		title + "\n" +
			bookRow + "\n" +
			shelfPlank.Render(strings.Repeat(" ", lipgloss.Width(bookRow))),
	)
}

