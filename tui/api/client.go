package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
    "tui/types"
)

type Client struct {
    BaseURL    string
    Token      string
    HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
    return &Client{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *Client) SetToken(token string) {
    c.Token = token
}

func (c *Client) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
    var reqBody io.Reader

    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        reqBody = bytes.NewBuffer(jsonData)
    }

    req, err := http.NewRequest(method, c.BaseURL+endpoint, reqBody)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    if c.Token != "" {
        req.Header.Set("Authorization", "Bearer "+c.Token)
    }

    return c.HTTPClient.Do(req)
}

// Auth endpoints
func (c *Client) Login(username, password string) (string, error) {
    data := map[string]string{
        "username": username,
        "password": password,
    }

    resp, err := c.doRequest("POST", "/auth/login", data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result struct {
        Token string `json:"token"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }

    if resp.StatusCode != 200 {
        return "", fmt.Errorf("login failed")
    }

    c.Token = result.Token
    return result.Token, nil
}

func (c *Client) GetCurrentUser() (types.User, error) {
    resp, err := c.doRequest("GET", "/me", nil)
    if err != nil {
        return types.User{}, err
    }
    defer resp.Body.Close()

    var user types.User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return types.User{}, err
    }

    return user, nil
}

// Books endpoints
func (c *Client) ListBooks() ([]types.Book, error) {
    resp, err := c.doRequest("GET", "/books", nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var books []types.Book
    if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
        return nil, err
    }

    return books, nil
}

func (c *Client) GetBook(bookID int) (types.Book, error) {
    resp, err := c.doRequest("GET", fmt.Sprintf("/books/%d", bookID), nil)
    if err != nil {
        return types.Book{}, err
    }
    defer resp.Body.Close()

    var book types.Book
    if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
        return types.Book{}, err
    }

    return book, nil
}

func (c *Client) AddReview(bookID int, text string, rating int) error {
    data := map[string]interface{}{
        "username": c.Token, // Using token as username for now
        "text":     text,
        "rating":   rating,
    }

    resp, err := c.doRequest("POST", fmt.Sprintf("/books/%d/reviews", bookID), data)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("failed to add review")
    }

    return nil
}

// Library endpoints
func (c *Client) GetUserLibraries(username string) ([]types.Library, error) {
    resp, err := c.doRequest("GET", fmt.Sprintf("/users/%s/libraries", username), nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var libraries []types.Library
    if err := json.NewDecoder(resp.Body).Decode(&libraries); err != nil {
        return nil, err
    }

    return libraries, nil
}

func (c *Client) CreateLibrary(name string) (int, error) {
    data := map[string]string{
        "username": c.Token,
        "name":     name,
    }

    resp, err := c.doRequest("POST", "/libraries", data)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    var result struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return 0, err
    }

    return result.ID, nil
}

func (c *Client) AddBookToLibrary(libraryID, bookID int, shelf string) error {
    data := map[string]interface{}{
        "username": c.Token,
        "book_id":  bookID,
        "shelf":    shelf,
    }

    resp, err := c.doRequest("POST", fmt.Sprintf("/libraries/%d/books", libraryID), data)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("failed to add book to library")
    }

    return nil
}

// Reading endpoints
func (c *Client) StartReading(bookID int) error {
    data := map[string]interface{}{
        "username": c.Token,
        "book_id":  bookID,
    }

    resp, err := c.doRequest("POST", "/reading/start", data)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("failed to start reading")
    }

    return nil
}

func (c *Client) TurnPage(bookID int, direction string, count int) error {
    data := map[string]interface{}{
        "username":  c.Token,
        "book_id":   bookID,
        "direction": direction,
        "count":     count,
    }

    resp, err := c.doRequest("POST", "/reading/turn", data)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("failed to turn page")
    }

    return nil
}

func (c *Client) GetActiveReading() ([]map[string]interface{}, error) {
    resp, err := c.doRequest("GET", fmt.Sprintf("/users/%s/reading", c.Token), nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var sessions []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&sessions); err != nil {
        return nil, err
    }

    return sessions, nil
}

// Friends endpoints
func (c *Client) AddFriend(friendUsername string) error {
    resp, err := c.doRequest("POST", fmt.Sprintf("/users/%s/friends/%s", c.Token, friendUsername), nil)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("failed to add friend")
    }

    return nil
}

func (c *Client) GetUser(username string) (types.User, error) {
    resp, err := c.doRequest("GET", fmt.Sprintf("/users/%s", username), nil)
    if err != nil {
        return types.User{}, err
    }
    defer resp.Body.Close()

    var user types.User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return types.User{}, err
    }

    return user, nil
}

// Recommendations endpoints
func (c *Client) RecommendBook(toUser string, bookID int, message string) error {
    data := map[string]interface{}{
        "from_user": c.Token,
        "to_user":   toUser,
        "book_id":   bookID,
        "message":   message,
    }

    resp, err := c.doRequest("POST", "/recommend", data)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("failed to send recommendation")
    }

    return nil
}

func (c *Client) GetRecommendations() ([]types.Recommendation, error) {
    resp, err := c.doRequest("GET", fmt.Sprintf("/users/%s/recommendations", c.Token), nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var recommendations []types.Recommendation
    if err := json.NewDecoder(resp.Body).Decode(&recommendations); err != nil {
        return nil, err
    }

    return recommendations, nil
}