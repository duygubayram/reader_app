CREATE TABLE IF NOT EXISTS recommendations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_user TEXT NOT NULL,
    to_user TEXT NOT NULL,
    book_id INTEGER NOT NULL,
    message TEXT,
    date TIMESTAMP NOT NULL,
    FOREIGN KEY (from_user) REFERENCES users(username) ON DELETE CASCADE,
    FOREIGN KEY (to_user) REFERENCES users(username) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);
