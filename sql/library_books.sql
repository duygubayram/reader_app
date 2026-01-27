CREATE TABLE IF NOT EXISTS library_books (
    library_id INTEGER NOT NULL,
    book_id INTEGER NOT NULL,
    shelf TEXT NOT NULL,
    PRIMARY KEY (library_id, book_id),
    FOREIGN KEY (library_id) REFERENCES libraries(id),
    FOREIGN KEY (book_id) REFERENCES books(id)
);
