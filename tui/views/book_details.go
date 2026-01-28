package views

import (
    "fmt"
    "strings"
    "github.com/charmbracelet/lipgloss"
    "tui/styles"
    "tui/types"
)

func RenderBookDetails(book types.Book, reviews []types.Review) string {
    header := styles.TitleStyle.Render(strings.ToUpper(book.Name))

    meta := []string{
        fmt.Sprintf("📖 Title: %s", book.Name),
        fmt.Sprintf("✍️  Author: %s", book.Author),
        fmt.Sprintf("📅 Year: %d", book.Year),
        fmt.Sprintf("📄 Pages: %d", book.Pages),
        fmt.Sprintf("⭐ Rating: %.1f/5", book.Rating),
        fmt.Sprintf("🌐 Language: %s", book.Language),
        fmt.Sprintf("🏢 Publisher: %s", book.Publisher),
    }

    // Reviews section
    reviewsSection := "\n📝 Reviews:\n"
    if len(reviews) == 0 {
        reviewsSection += "  No reviews yet\n"
    } else {
        for _, review := range reviews {
            reviewsSection += fmt.Sprintf("  %s: ⭐%d - %s\n", review.User, review.Rating, review.Text)
        }
    }

    actions := lipgloss.JoinHorizontal(
        lipgloss.Top,
        styles.ButtonStyle.Render("📖 Start Reading"),
        styles.ButtonStyle.Render("➕ Add to Library"),
        styles.ButtonStyle.Render("💬 Add Review"),
        styles.ButtonStyle.Render("⭐ Rate Book"),
    )

    content := lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        "\n",
        strings.Join(meta, "\n"),
        reviewsSection,
        "\n",
        actions,
    )

    return styles.CardStyle.Width(60).Render(content)
}