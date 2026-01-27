package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* =======================
   CONFIG
======================= */

const api = "http://localhost:8000"
const currentUser = "alice"

/* =======================
   STYLES
======================= */

var (
	windowStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	tabStyle = lipgloss.NewStyle().
			Padding(0, 2)

	activeTabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	contentStyle = lipgloss.NewStyle().
			Padding(1, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))
)

/* =======================
   TYPES
======================= */

type tab int

const (
	profileTab tab = iota
	librariesTab
	booksTab
	readingTab
	recommendTab
)

/* ---------- API Models ---------- */

type User struct {
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Friends     []string `json:"friends"`
	Libraries   []string `json:"libraries"`
}

type Book struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Author    string   `json:"author"`
	Year      int      `json:"year"`
	Language  string   `json:"language"`
	Publisher string   `json:"publisher"`
	Pages     int      `json:"pages"`
	AvgRating *float64 `json:"avg_rating"`
}

type ReadingSession struct {
	User        string `json:"user"`
	BookID      int    `json:"book_id"`
	CurrentPage int    `json:"current_page"`
}

/* ---------- Bubble Tea Model ---------- */

type model struct {
	activeTab tab

	user    User
	books   []Book
	reading []ReadingSession
	recs    []map[string]any

	cursor  int
	err     error
	loading bool
}

/* =======================
   INIT
======================= */

func initialModel() model {
	return model{
		activeTab: profileTab,
	}
}

func (m model) Init() tea.Cmd {
	return loadUser()
}

/* =======================
   API COMMANDS
======================= */

func loadUser() tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get(api + "/users/" + currentUser)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var u User
		json.NewDecoder(resp.Body).Decode(&u)
		return u
	}
}

func loadBooks() tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get(api + "/books")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var b []Book
		json.NewDecoder(resp.Body).Decode(&b)
		return b
	}
}

func loadReading() tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get(api + "/users/" + currentUser + "/reading")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var r []ReadingSession
		json.NewDecoder(resp.Body).Decode(&r)
		return r
	}
}

func startReading(bookID int) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]any{
			"username": currentUser,
			"book_id":  bookID,
		})
		http.Post(api+"/reading/start", "application/json", bytes.NewBuffer(body))
		return loadReading()()
	}
}

func turnPage(bookID int, dir string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]any{
			"username":  currentUser,
			"book_id":   bookID,
			"direction": dir,
		})
		http.Post(api+"/reading/turn", "application/json", bytes.NewBuffer(body))
		return loadReading()()
	}
}

/* =======================
   UPDATE
======================= */

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case User:
		m.user = msg
		return m, nil

	case []Book:
		m.books = msg
		return m, nil

	case []ReadingSession:
		m.reading = msg
		return m, nil

	case error:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "left":
			if m.activeTab > 0 {
				m.activeTab--
			}

		case "right":
			if m.activeTab < recommendTab {
				m.activeTab++
			}

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			m.cursor++

		case "enter":
			if m.activeTab == booksTab && len(m.books) > 0 {
				return m, startReading(m.books[m.cursor].ID)
			}

		case "h":
			if m.activeTab == readingTab && len(m.reading) > 0 {
				return m, turnPage(m.reading[m.cursor].BookID, "back")
			}

		case "l":
			if m.activeTab == readingTab && len(m.reading) > 0 {
				return m, turnPage(m.reading[m.cursor].BookID, "forward")
			}
		}
	}

	/* Tab side effects */
	switch m.activeTab {
	case booksTab:
		if m.books == nil {
			return m, loadBooks()
		}
	case readingTab:
		return m, loadReading()
	}

	return m, nil
}

/* =======================
   VIEW
======================= */

func (m model) View() string {
	tabs := m.renderTabs()
	content := m.renderContent()
	ui := lipgloss.JoinVertical(lipgloss.Left, tabs, content)
	return windowStyle.Render(ui)
}

func (m model) renderTabs() string {
	labels := []string{"Profile", "Libraries", "Books", "Reading", "Inbox"}
	var out string
	for i, label := range labels {
		if tab(i) == m.activeTab {
			out += activeTabStyle.Render(label)
		} else {
			out += tabStyle.Render(label)
		}
	}
	return out
}

func (m model) renderContent() string {
	if m.err != nil {
		return errorStyle.Render(m.err.Error())
	}

	switch m.activeTab {

	case profileTab:
		return fmt.Sprintf(
			"Profile\n\nUser: %s\nDisplay: %s\nFriends: %d\nLibraries: %d",
			m.user.Username,
			m.user.DisplayName,
			len(m.user.Friends),
			len(m.user.Libraries),
		)

	case librariesTab:
		out := "Libraries\n\n"
		for _, l := range m.user.Libraries {
			out += "• " + l + "\n"
		}
		return out

	case booksTab:
		out := "Books (Enter = start reading)\n\n"
		for i, b := range m.books {
			cursor := " "
			if i == m.cursor {
				cursor = ">"
			}
			out += fmt.Sprintf("%s %s — %s (%d)\n", cursor, b.Name, b.Author, b.Year)
		}
		return out

	case readingTab:
		out := "Reading (h/l = turn pages)\n\n"
		for _, r := range m.reading {
			out += fmt.Sprintf("Book %d → page %d\n", r.BookID, r.CurrentPage)
		}
		return out

	case recommendTab:
		return "Recommendations\n\n(coming from /users/{u}/recommendations)"

	}

	return ""
}

/* =======================
   MAIN
======================= */

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
