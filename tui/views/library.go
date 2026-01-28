package views

import (
    "tui/types"
    "github.com/charmbracelet/lipgloss"
)

func RenderLibrary(shelves map[string][]types.Book, selectedShelf, selectedBook int) string {
    shelfNames := []string{"to_read", "currently_reading", "read"}

    var renderedShelves []string

    for i, shelfName := range shelfNames {
        books := shelves[shelfName]
        isSelected := i == selectedShelf

        shelfContent := RenderShelf(shelfName, books, selectedBook, isSelected)
        if isSelected {
            shelfContent = lipgloss.NewStyle().
                Padding(1, 2).
                Border(lipgloss.RoundedBorder()).
                BorderForeground(lipgloss.Color("#2563EB")).
                Background(lipgloss.Color("#374151")).
                Render(shelfContent)
        } else {
            shelfContent = lipgloss.NewStyle().
                Padding(1, 2).
                Border(lipgloss.RoundedBorder()).
                BorderForeground(lipgloss.Color("#4B5563")).
                Background(lipgloss.Color("#374151")).
                Render(shelfContent)
        }

        renderedShelves = append(renderedShelves, shelfContent)
    }

    return lipgloss.JoinVertical(lipgloss.Top, renderedShelves...)
}