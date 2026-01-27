CREATE TABLE IF NOT EXISTS reading_sessions (
    user TEXT NOT NULL,
    book_id INTEGER NOT NULL,
    current_page INTEGER NOT NULL,
    started_at TIMESTAMP NOT NULL,
    last_read_at TIMESTAMP,
    PRIMARY KEY (user, book_id),
    FOREIGN KEY (user) REFERENCES users(username) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);
