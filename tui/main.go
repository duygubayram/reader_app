package main

import (
    "fmt"
    "os"
    tea "github.com/charmbracelet/bubbletea"
    "tui/app"
)

func main() {
    // Set API URL from environment or use default
    apiURL := "http://localhost:8000"
    if url := os.Getenv("BOOKTRACKER_API_URL"); url != "" {
        apiURL = url
    }

    m := app.NewModel(apiURL)

    p := tea.NewProgram(m,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
        // tea.WithFPS(60),
    )

    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}