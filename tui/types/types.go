package types

// View types
type View int

const (
    ViewLogin View = iota
    ViewLibrary
    ViewBookDetails
    ViewReading
    ViewProfile
    ViewFriends
    ViewRecommendations
)

// Model types
type NavItem struct {
    ID    string
    Label string
    View  View
}

type LoginForm struct {
    Username string
    Password string
    Focused  string // "username" or "password"
}

func (f *LoginForm) UpdateUsername(msg interface{}) {
    // Implement in your actual code
}

func (f *LoginForm) UpdatePassword(msg interface{}) {
    // Implement in your actual code
}

func (f LoginForm) RenderUsername() string {
    if f.Focused == "username" {
        return f.Username + "█"
    }
    return f.Username
}

func (f LoginForm) RenderPassword() string {
    if f.Focused == "password" {
        if f.Password == "" {
            return "█"
        }
        return "•" + "█"
    }
    return "••••••"
}

type SearchBar struct {
    Active  bool
    Query   string
    Results []Book
}

type BookList struct {
    Books    []Book
    Selected int
    Page     int
    PageSize int
}

type ShelfView struct {
    Shelves       map[string][]Book
    SelectedShelf int
    SelectedBook  int
}

type ReadingView struct {
    Book        Book
    CurrentPage int
    Progress    float64
}

type ProfileView struct {
    User        User
    Stats       UserStats
    RecentBooks []Book
}

// Data types
type Book struct {
    ID        int
    Name      string
    Author    string
    Year      int
    Pages     int
    Rating    float64
    Status    string // "to_read", "currently_reading", "read"
    Language  string
    Publisher string
}

type User struct {
    Username    string
    DisplayName string
    JoinedDate  string
    Friends     []string
    Libraries   []Library
}

type Library struct {
    ID    int
    Name  string
    Books map[string][]int
}

type UserStats struct {
    TotalBooks   int
    BooksReading int
    BooksRead    int
    Friends      int
    AvgRating    float64
}

type Friend struct {
    Username    string
    DisplayName string
    Online      bool
    Reading     string
}

type Recommendation struct {
    From    string
    Book    string
    Message string
    Date    string
}

type Review struct {
    User   string
    Rating int
    Text   string
    Likes  int
}

type Activity struct {
    Type      string // "read", "reviewed", "added_book", "added_friend"
    BookID    int
    BookTitle string
    Date      string
}

type LibraryData struct {
    TotalBooks       int
    Shelves          map[string][]Book
    RecentAdds       []Book
    CurrentlyReading []Book
}

type BookData struct {
    Book    Book
    Reviews []Review
    Similar []Book
}

type ProfileData struct {
    User     User
    Stats    UserStats
    Activity []Activity
}

// Message types
type LoginSuccessMsg struct {
    Username string
    Token    string
    User     User
}

type LoginErrorMsg struct {
    Message string
}

type LoadLibraryMsg struct {
    Books   []Book
    Shelves map[string][]Book
}

type LoadUserMsg struct {
    User User
}

type LoadFriendsMsg struct {
    Friends []Friend
}

type LoadRecommendationsMsg struct {
    Recommendations []Recommendation
}

type LoadReadingSessionsMsg struct {
    Sessions []map[string]interface{}
}

type ErrorMsg struct {
    Message string
}

type SwitchToReadingMsg struct {
    BookID int
}

type ClearErrorMsg struct{}

type ApiResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}