package views

import (
    "strings"
    "github.com/charmbracelet/lipgloss"
    "tui/styles"
    "tui/types"
)

func RenderShelf(name string, books []types.Book, selectedBook int, isSelected bool) string {
    var renderedBooks []string

    for i, book := range books {
        bookSelected := isSelected && i == selectedBook
        renderedBooks = append(renderedBooks, RenderBook(book, bookSelected))
    }

    if len(renderedBooks) == 0 {
        renderedBooks = []string{styles.BookStyle.Render("Empty\nShelf")}
    }

    bookRow := lipgloss.JoinHorizontal(
        lipgloss.Top,
        renderedBooks...,
    )

    title := styles.ShelfPlankStyle.Render(" " + strings.ReplaceAll(name, "_", " ") + " ")

    return lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        "\n",
        bookRow,
    )
}

func RenderBook(book types.Book, selected bool) string {
    title := book.Name
    if len(title) > 8 {
        title = title[:8] + ".."
    }

    lines := []string{
        "",
        title,
        "─────",
        book.Author,
    }

    content := strings.Join(lines, "\n")

    style := styles.BookStyle
    if selected {
        style = styles.BookSelectedStyle
    }

    return style.Render(content)
}