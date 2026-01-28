package app

import (
    tea "github.com/charmbracelet/bubbletea"
    "tui/api"
    "tui/types"
    "time"
)

type Model struct {
    width  int
    height int

    // API client
    api *api.Client

    // App state
    currentView types.View
    loading     bool
    errorMsg    string

    // Navigation
    navItems    []types.NavItem
    selectedNav int

    // Authentication
    username string
    token    string
    loggedIn bool

    // Data
    libraryData      types.LibraryData
    bookData         types.BookData
    profileData      types.ProfileData
    friendsData      []types.Friend
    recommendations  []types.Recommendation
    activeReading    []map[string]interface{}

    // UI Components
    loginForm   types.LoginForm
    searchBar   types.SearchBar
    bookList    types.BookList
    shelfView   types.ShelfView
    readingView types.ReadingView
    profileView types.ProfileView

    // Selected book for details view
    selectedBookID int
}

func NewModel(apiURL string) Model {
    apiClient := api.NewClient(apiURL)

    navItems := []types.NavItem{
        {ID: "library", Label: "📚 My Library", View: types.ViewLibrary},
        {ID: "discover", Label: "🔍 Discover", View: types.ViewLibrary},
        {ID: "reading", Label: "📖 Reading", View: types.ViewReading},
        {ID: "friends", Label: "👥 Friends", View: types.ViewFriends},
        {ID: "recommendations", Label: "💡 Recommendations", View: types.ViewRecommendations},
        {ID: "profile", Label: "👤 Profile", View: types.ViewProfile},
    }

    return Model{
        api:          apiClient,
        currentView:  types.ViewLogin,
        navItems:     navItems,
        selectedNav:  0,
        loginForm: types.LoginForm{
            Username: "",
            Password: "",
            Focused:  "username",
        },
        libraryData: types.LibraryData{
            Shelves: make(map[string][]types.Book),
        },
        shelfView: types.ShelfView{
            Shelves: make(map[string][]types.Book),
        },
    }
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    // Handle API responses first
    switch msg := msg.(type) {
    case types.LoginSuccessMsg:
        m.loading = false
        m.loggedIn = true
        m.username = msg.Username
        m.token = msg.Token
        m.currentView = types.ViewLibrary
        m.profileData.User = msg.User
        m.api.SetToken(msg.Token)

        // Load initial data
        return m, tea.Batch(
            m.loadLibraryData(),
            m.loadProfileData(),
        )

    case types.LoginErrorMsg:
        m.loading = false
        m.errorMsg = msg.Message
        return m, m.clearErrorAfter(3)

    case types.LoadLibraryMsg:
        m.loading = false
        m.libraryData.Shelves = msg.Shelves
        m.shelfView.Shelves = msg.Shelves
        return m, nil

    case types.LoadUserMsg:
        m.loading = false
        m.profileData.User = msg.User
        return m, nil

    case types.ErrorMsg:
        m.loading = false
        m.errorMsg = msg.Message
        return m, m.clearErrorAfter(3)

    case types.ClearErrorMsg:
        m.errorMsg = ""
        return m, nil
    }

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil

    case tea.KeyMsg:
        cmd = m.handleKeyPress(msg)
    }

    // Delegate to specific view handlers
    switch m.currentView {
    case types.ViewLogin:
        return m.updateLogin(msg)
    case types.ViewLibrary:
        return m.updateLibrary(msg)
    case types.ViewBookDetails:
        return m.updateBookDetails(msg)
    case types.ViewReading:
        return m.updateReading(msg)
    case types.ViewProfile:
        return m.updateProfile(msg)
    default:
        return m, cmd
    }
}

func (m Model) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
    switch msg.String() {
    case "ctrl+c", "q":
        return tea.Quit
    case "tab":
        if m.currentView == types.ViewLogin {
            m.loginForm.Focused = "password"
        }
    case "shift+tab":
        if m.currentView == types.ViewLogin {
            m.loginForm.Focused = "username"
        }
    }
    return nil
}

func (m Model) View() string {
    if m.loading {
        return m.renderLoading()
    }

    if m.errorMsg != "" {
        return m.renderError()
    }

    // Render main layout with header, content, and footer
    return m.renderLayout()
}

// Helper to clear error after seconds
func (m Model) clearErrorAfter(seconds int) tea.Cmd {
    return tea.Tick(time.Second*time.Duration(seconds), func(t time.Time) tea.Msg {
        return types.ClearErrorMsg{}
    })
}