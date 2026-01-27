package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* ================= CONFIG ================= */

const api = "http://localhost:8000"

/* ================= STYLES ================= */

var (
	window = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(80)

	title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)
)

/* ================= TYPES ================= */

type screen int
type tab int

const (
	screenWelcome screen = iota
	screenLogin
	screenRegister
	screenMain
)

const (
	profileTab tab = iota
	booksTab
	readingTab
	inboxTab
)

/* ================= API MODELS ================= */

type User struct {
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Friends     []string `json:"friends"`
	Libraries   []string `json:"libraries"`
}

type Book struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

type ReadingSession struct {
	BookID      int `json:"book_id"`
	CurrentPage int `json:"current_page"`
}

type Recommendation struct {
	From    string `json:"from"`
	Book    string `json:"book"`
	Message string `json:"message"`
}

/* ================= MODEL ================= */

type model struct {
	screen screen
	tab    tab
	cursor int

	username string
	password string
	display  string

	user    User
	books   []Book
	reading []ReadingSession
	recs    []Recommendation

	input string
	err   string
}

/* ================= INIT ================= */

func initialModel() model {
	return model{screen: screenWelcome}
}

func (m model) Init() tea.Cmd {
	return nil
}

/* ================= API ================= */

func login(username, password string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{
			"username": username,
			"password": password,
		})
		resp, err := http.Post(api+"/auth/login", "application/json", bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("login failed")
		}
		return username
	}
}

func register(username, display, password string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{
			"username":     username,
			"display_name": display,
			"password":     password,
		})
		resp, err := http.Post(api+"/register", "application/json", bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("register failed")
		}
		return username
	}
}

func loadUser(username string) tea.Cmd {
	return func() tea.Msg {
		resp, _ := http.Get(api + "/users/" + username)
		defer resp.Body.Close()
		var u User
		json.NewDecoder(resp.Body).Decode(&u)
		return u
	}

}

func loadBooks() tea.Cmd {
	return func() tea.Msg {
		resp, _ := http.Get(api + "/books")
		defer resp.Body.Close()
		var b []Book
		json.NewDecoder(resp.Body).Decode(&b)
		return b
	}
}

func loadReading(username string) tea.Cmd {
	return func() tea.Msg {
		resp, _ := http.Get(api + "/users/" + username + "/reading")
		defer resp.Body.Close()
		var r []ReadingSession
		json.NewDecoder(resp.Body).Decode(&r)
		return r
	}
}

func loadInbox(username string) tea.Cmd {
	return func() tea.Msg {
		resp, _ := http.Get(api + "/users/" + username + "/recommendations")
		defer resp.Body.Close()
		var r []Recommendation
		json.NewDecoder(resp.Body).Decode(&r)
		return r
	}
}

func startReading(username string, bookID int) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]any{
			"username": username,
			"book_id":  bookID,
		})
		http.Post(api+"/reading/start", "application/json", bytes.NewBuffer(body))
		return loadReading(username)()
	}
}

func turnPage(username string, bookID int, dir string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]any{
			"username":  username,
			"book_id":   bookID,
			"direction": dir,
		})
		http.Post(api+"/reading/turn", "application/json", bytes.NewBuffer(body))
		return loadReading(username)()
	}
}

/* ================= UPDATE ================= */

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case string: // login success
		m.username = msg
		m.screen = screenMain
		return m, tea.Batch(
			loadUser(m.username),
			loadBooks(),
			loadReading(m.username),
			loadInbox(m.username),
		)

	case User:
		m.user = msg

	case []Book:
		m.books = msg

	case []ReadingSession:
		m.reading = msg

	case []Recommendation:
		m.recs = msg

	case error:
		m.err = msg.Error()

	case tea.KeyMsg:
		k := msg.String()

		if k == "ctrl+c" || k == "q" {
			return m, tea.Quit
		}

		switch m.screen {

		case screenWelcome:
			if k == "l" {
				m.screen = screenLogin
			}
			if k == "n" {
				m.screen = screenRegister
			}

		case screenLogin:
			if k == "enter" {
				parts := strings.Split(m.input, ",")
				if len(parts) == 2 {
					return m, login(parts[0], parts[1])
				}
			}
			if k == "backspace" && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			} else {
				m.input += k
			}

		case screenRegister:
			if k == "enter" {
				parts := strings.Split(m.input, ",")
				if len(parts) == 3 {
					return m, register(parts[0], parts[1], parts[2])
				}
			}
			if k == "backspace" && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			} else {
				m.input += k
			}

		case screenMain:
			switch k {
			case "left":
				if m.tab > 0 {
					m.tab--
				}
			case "right":
				if m.tab < inboxTab {
					m.tab++
				}
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				m.cursor++
			case "enter":
				if m.tab == booksTab && len(m.books) > 0 {
					return m, startReading(m.username, m.books[m.cursor].ID)
				}
			case "h":
				if m.tab == readingTab {
					return m, turnPage(m.username, m.reading[m.cursor].BookID, "back")
				}
			case "l":
				if m.tab == readingTab {
					return m, turnPage(m.username, m.reading[m.cursor].BookID, "forward")
				}
			}
		}
	}

	return m, nil
}

/* ================= VIEW ================= */

func (m model) View() string {

	switch m.screen {

	case screenWelcome:
		return window.Render(title.Render("Reader App") + "\n\n[l] Login\n[n] Register\n[q] Quit")

	case screenLogin:
		return window.Render("Login\n\nusername,password\n\n> " + m.input)

	case screenRegister:
		return window.Render("Register\n\nusername,display,password\n\n> " + m.input)

	case screenMain:
		tabs := []string{"Profile", "Books", "Reading", "Inbox"}
		bar := ""
		for i, t := range tabs {
			if tab(i) == m.tab {
				bar += cursorStyle.Render(" " + t + " ")
			} else {
				bar += " " + t + " "
			}
		}

		return window.Render(
			title.Render(bar) + "\n\n" +
				m.renderContent() +
				"\n\n" +
				help.Render("← → tabs | ↑ ↓ cursor | enter select | h/l pages | q quit"),
		)
	}

	return ""
}

func (m model) renderContent() string {

	switch m.tab {

	case profileTab:
		return fmt.Sprintf(
			"User: %s\nDisplay: %s\nFriends: %d\nLibraries: %d",
			m.user.Username,
			m.user.DisplayName,
			len(m.user.Friends),
			len(m.user.Libraries),
		)

	case booksTab:
		out := "Books:\n\n"
		for i, b := range m.books {
			prefix := " "
			if i == m.cursor {
				prefix = ">"
			}
			out += fmt.Sprintf("%s %s — %s\n", prefix, b.Name, b.Author)
		}
		return out

	case readingTab:
		out := "Reading:\n\n"
		for _, r := range m.reading {
			out += fmt.Sprintf("Book %d → page %d\n", r.BookID, r.CurrentPage)
		}
		return out

	case inboxTab:
		out := "Inbox:\n\n"
		for _, r := range m.recs {
			out += fmt.Sprintf("From %s → %s\n%s\n\n",
				r.From, r.Book, r.Message)
		}
		return out
	}

	return ""
}

/* ================= MAIN ================= */

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
