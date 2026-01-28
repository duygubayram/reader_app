package app

import (
    "strings"
    "github.com/charmbracelet/lipgloss"
    "tui/views"
    "tui/types"
)

func (m Model) renderLayout() string {
    header := m.renderHeader()
    content := m.renderContent()
    footer := m.renderFooter()

    layout := lipgloss.JoinVertical(
        lipgloss.Top,
        header,
        content,
        footer,
    )

    // Apply styling with terminal dimensions
    return lipgloss.NewStyle().
        Width(m.width).
        Height(m.height).
        Render(layout)
}

func (m Model) renderHeader() string {
    if m.currentView == types.ViewLogin {
        return lipgloss.NewStyle().
            Background(lipgloss.Color("#2563EB")).
            Foreground(lipgloss.Color("#F3F4F6")).
            Padding(0, 2).
            Height(3).
            Bold(true).
            Render("📚 Book Tracker")
    }

    greeting := ""
    if m.username != "" {
        greeting = " | Welcome, " + m.username
    }

    left := lipgloss.NewStyle().
        Bold(true).
        Padding(0, 1).
        Render("📚 Book Tracker" + greeting)

    right := lipgloss.NewStyle().
        Faint(true).
        Render("📅 Today's Read")

    return lipgloss.NewStyle().
        Background(lipgloss.Color("#2563EB")).
        Foreground(lipgloss.Color("#F3F4F6")).
        Padding(0, 2).
        Height(3).
        Bold(true).
        Render(
            lipgloss.JoinHorizontal(lipgloss.Center, left, right),
        )
}

func (m Model) renderContent() string {
    switch m.currentView {
    case types.ViewLogin:
        return m.renderLoginView()
    case types.ViewLibrary:
        return m.renderLibraryView()
    case types.ViewBookDetails:
        return m.renderBookDetailsView()
    case types.ViewReading:
        return m.renderReadingView()
    case types.ViewProfile:
        return m.renderProfileView()
    default:
        return "Coming soon..."
    }
}

func (m Model) renderFooter() string {
    helpText := ""
    switch m.currentView {
    case types.ViewLogin:
        helpText = "↑↓: Navigate | Tab: Next field | Enter: Login | Q: Quit"
    case types.ViewLibrary:
        helpText = "←→: Move shelf | ↑↓: Move book | Enter: Select | N: New book | S: Search | Q: Quit"
    case types.ViewBookDetails:
        helpText = "R: Start reading | A: Add to library | F: Add friend | Esc: Back | Q: Quit"
    }

    status := ""
    if m.loggedIn {
        status = "🟢 Online"
    } else {
        status = "🔴 Offline"
    }

    left := lipgloss.NewStyle().
        Faint(true).
        Render(helpText)

    right := lipgloss.NewStyle().
        Bold(true).
        Render(status)

    return lipgloss.NewStyle().
        Background(lipgloss.Color("#1F2937")).
        Foreground(lipgloss.Color("#F3F4F6")).
        Padding(0, 2).
        Height(2).
        Render(
            lipgloss.JoinHorizontal(lipgloss.Center, left, right),
        )
}

func (m Model) renderLoginView() string {
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#2563EB")).
        MarginBottom(1).
        Render("Welcome to Book Tracker")

    subtitle := lipgloss.NewStyle().
        Faint(true).
        MarginBottom(2).
        Render("Your personal digital library")

    form := lipgloss.NewStyle().
        Width(40).
        Padding(2, 3).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#2563EB")).
        Render(
            lipgloss.JoinVertical(lipgloss.Left,
                lipgloss.NewStyle().Bold(true).MarginBottom(1).Render("Username:"),
                m.loginForm.RenderUsername(),
                "\n",
                lipgloss.NewStyle().Bold(true).MarginBottom(1).Render("Password:"),
                m.loginForm.RenderPassword(),
                "\n",
                lipgloss.NewStyle().
                    Background(lipgloss.Color("#2563EB")).
                    Foreground(lipgloss.Color("#F3F4F6")).
                    Padding(0, 3).
                    Bold(true).
                    Render("Login"),
            ),
        )

    return lipgloss.JoinVertical(
        lipgloss.Center,
        "\n\n",
        title,
        subtitle,
        "\n\n",
        form,
    )
}

func (m Model) renderLibraryView() string {
    navBar := m.renderNavBar()
    mainContent := ""

    if m.searchBar.Active {
        mainContent = m.renderSearchView()
    } else {
        mainContent = m.renderShelfView()
    }

    sidebar := m.renderSidebar()

    content := lipgloss.JoinHorizontal(
        lipgloss.Top,
        sidebar,
        mainContent,
    )

    return lipgloss.JoinVertical(
        lipgloss.Top,
        navBar,
        content,
    )
}

func (m Model) renderNavBar() string {
    var navItems []string
    for i, item := range m.navItems {
        style := lipgloss.NewStyle().
            Padding(0, 2).
            Foreground(lipgloss.Color("#F3F4F6"))

        if i == m.selectedNav {
            style = style.
                Background(lipgloss.Color("#7C3AED")).
                Bold(true)
        }
        navItems = append(navItems, style.Render(item.Label))
    }

    return lipgloss.NewStyle().
        Background(lipgloss.Color("#374151")).
        Padding(0, 2).
        Height(2).
        Render(
            lipgloss.JoinHorizontal(lipgloss.Left, navItems...),
        )
}

func (m Model) renderSidebar() string {
    stats := []string{
        "📊 Your Stats",
        "─────────────",
        "📚 Books: 42",
        "📖 Reading: 3",
        "✅ Read: 12",
        "👥 Friends: 8",
        "⭐ Avg Rating: 4.2",
    }

    return lipgloss.NewStyle().
        Width(25).
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#2563EB")).
        Render(
            strings.Join(stats, "\n"),
        )
}

func (m Model) renderShelfView() string {
    // Use the views package to render the shelf
    if m.shelfView.Shelves == nil || len(m.shelfView.Shelves) == 0 {
        return "No books in library yet."
    }
    return views.RenderLibrary(m.shelfView.Shelves, m.shelfView.SelectedShelf, m.shelfView.SelectedBook)
}

func (m Model) renderBookDetailsView() string {
    // Fetch book details if not already loaded
    if m.selectedBookID > 0 {
        // We would typically have loaded the book data, but for now, let's use a placeholder.
        // In a real app, you would have a method to load the book data and reviews.
        book := types.Book{
            ID:       m.selectedBookID,
            Name:     "Sample Book",
            Author:   "Unknown",
            Year:     2023,
            Pages:    300,
            Rating:   4.5,
            Language: "English",
            Publisher: "Sample Publisher",
        }
        return views.RenderBookDetails(book, []types.Review{})
    }
    return "No book selected"
}

func (m Model) renderReadingView() string {
    return "Reading view not implemented"
}

func (m Model) renderProfileView() string {
    return "Profile view not implemented"
}

func (m Model) renderSearchView() string {
    return "Search view not implemented"
}

func (m Model) renderLoading() string {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("#F59E0B")).
        Bold(true).
        Render("Loading...")
}

func (m Model) renderError() string {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("#EF4444")).
        Bold(true).
        Render("Error: " + m.errorMsg)
}