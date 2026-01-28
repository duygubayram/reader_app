package app

import (
    tea "github.com/charmbracelet/bubbletea"
    "tui/types"
)

func (m Model) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            // Attempt login
            m.loading = true
            return m, m.attemptLogin()
        case "up", "down":
            if m.loginForm.Focused == "username" {
                m.loginForm.Focused = "password"
            } else {
                m.loginForm.Focused = "username"
            }
        }
    }

    // Update form fields
    if m.loginForm.Focused == "username" {
        m.loginForm.UpdateUsername(msg)
    } else {
        m.loginForm.UpdatePassword(msg)
    }

    return m, nil
}

func (m Model) updateLibrary(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "left":
            if m.selectedNav > 0 {
                m.selectedNav--
            }
        case "right":
            if m.selectedNav < len(m.navItems)-1 {
                m.selectedNav++
            }
        case "enter":
            if m.selectedNav < len(m.navItems) {
                // Convert NavItem.View (types.View) to app.View
                navView := m.navItems[m.selectedNav].View
                m.currentView = navView

                // Load data for the selected view
                switch m.currentView {
                case types.ViewLibrary:
                    return m, m.loadLibraryData()
                case types.ViewProfile:
                    return m, m.loadProfileData()
                case types.ViewFriends:
                    return m, m.loadFriendsData()
                case types.ViewRecommendations:
                    return m, m.loadRecommendations()
                case types.ViewReading:
                    return m, m.loadReadingSessions()
                }
            }
        case "s":
            m.searchBar.Active = !m.searchBar.Active
        case "r":
            // Refresh data
            return m, m.refreshData()
        case "esc":
            m.currentView = types.ViewLibrary
        }
    }

    return m, nil
}

func (m Model) updateBookDetails(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc", "backspace":
            m.currentView = types.ViewLibrary
            m.selectedBookID = 0
        case "r":
            // Start reading the book
            if m.selectedBookID > 0 {
                return m, m.startReading(m.selectedBookID)
            }
        case "a":
            // Add review
            if m.selectedBookID > 0 {
                // You could show a form for adding review here
                m.errorMsg = "Review feature coming soon!"
                return m, m.clearErrorAfter(2)
            }
        }
    }
    return m, nil
}

func (m Model) updateReading(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle reading view updates
    return m, nil
}

func (m Model) updateProfile(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle profile view updates
    return m, nil
}

func (m Model) attemptLogin() tea.Cmd {
    return func() tea.Msg {
        // Call API
        token, err := m.api.Login(m.loginForm.Username, m.loginForm.Password)
        if err != nil {
            return types.LoginErrorMsg{Message: err.Error()}
        }

        // Get user data
        user, err := m.api.GetCurrentUser()
        if err != nil {
            return types.LoginErrorMsg{Message: err.Error()}
        }

        return types.LoginSuccessMsg{
            Username: m.loginForm.Username,
            Token:    token,
            User:     user,
        }
    }
}

func (m Model) loadLibraryData() tea.Cmd {
    return func() tea.Msg {
        // Load books
        books, err := m.api.ListBooks()
        if err != nil {
            return types.ErrorMsg{Message: err.Error()}
        }

        // Organize by shelves
        shelves := make(map[string][]types.Book)
        shelves["to_read"] = []types.Book{}
        shelves["currently_reading"] = []types.Book{}
        shelves["read"] = []types.Book{}

        // Load user's libraries
        libraries, err := m.api.GetUserLibraries(m.username)
        if err == nil && len(libraries) > 0 {
            // Use first library's organization
            for shelf, bookIDs := range libraries[0].Books {
                for _, id := range bookIDs {
                    for _, book := range books {
                        if book.ID == id {
                            shelves[shelf] = append(shelves[shelf], book)
                            break
                        }
                    }
                }
            }
        }

        return types.LoadLibraryMsg{
            Books:   books,
            Shelves: shelves,
        }
    }
}

func (m Model) loadProfileData() tea.Cmd {
    return func() tea.Msg {
        user, err := m.api.GetUser(m.username)
        if err != nil {
            return types.ErrorMsg{Message: err.Error()}
        }

        return types.LoadUserMsg{User: user}
    }
}

func (m Model) loadFriendsData() tea.Cmd {
    return func() tea.Msg {
        // Get user data which includes friends
        user, err := m.api.GetUser(m.username)
        if err != nil {
            return types.ErrorMsg{Message: err.Error()}
        }

        // Convert to friends list
        var friends []types.Friend
        for _, friendUsername := range user.Friends {
            friendUser, err := m.api.GetUser(friendUsername)
            if err == nil {
                friends = append(friends, types.Friend{
                    Username:    friendUser.Username,
                    DisplayName: friendUser.DisplayName,
                    Online:      false, // You'd need an online status endpoint
                })
            }
        }

        return types.LoadFriendsMsg{Friends: friends}
    }
}

func (m Model) loadRecommendations() tea.Cmd {
    return func() tea.Msg {
        recommendations, err := m.api.GetRecommendations()
        if err != nil {
            return types.ErrorMsg{Message: err.Error()}
        }

        return types.LoadRecommendationsMsg{Recommendations: recommendations}
    }
}

func (m Model) loadReadingSessions() tea.Cmd {
    return func() tea.Msg {
        sessions, err := m.api.GetActiveReading()
        if err != nil {
            return types.ErrorMsg{Message: err.Error()}
        }

        return types.LoadReadingSessionsMsg{Sessions: sessions}
    }
}

func (m Model) startReading(bookID int) tea.Cmd {
    return func() tea.Msg {
        err := m.api.StartReading(bookID)
        if err != nil {
            return types.ErrorMsg{Message: err.Error()}
        }

        // Switch to reading view
        return types.SwitchToReadingMsg{BookID: bookID}
    }
}

func (m Model) refreshData() tea.Cmd {
    return tea.Batch(
        m.loadLibraryData(),
        m.loadProfileData(),
    )
}