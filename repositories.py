import sqlite3
from main import User, Book

class Database:
    def __init__(self, path="app.db"): # add sql file path
        self.conn = sqlite3.connect(path, check_same_thread=False)
        self.conn.row_factory = sqlite3.Row # dict-like rows
        self.conn.execute("PRAGMA foreign_keys = ON")

    def execute(self, query, params=()):
        cur = self.conn.execute(query, params)
        self.conn.commit()
        return cur

    def fetchone(self, query, params=()):
        cur = self.conn.execute(query, params)
        return cur.fetchone()

    def fetchall(self, query, params=()):
        cur = self.conn.execute(query, params)
        return cur.fetchall()

def init_db(db: Database):
    for file in [
        "sql/books.sql",
        "sql/users.sql",
        "sql/friends.sql",
        "sql/libraries.sql",
        "sql/library_books.sql",
        "sql/reading_sessions.sql",
        "sql/reviews.sql",
        "sql/recommendations.sql"
    ]:
        with open(file, "r", encoding="utf-8") as f:
            sql = f.read()
        db.conn.executescript(sql)
    db.conn.commit()
    #db.execute("""
    # CREATE TABLE IF NOT EXISTS users (
    #     username TEXT PRIMARY KEY,
    #     display_name TEXT NOT NULL,
    #     created_at TIMESTAMP NOT NULL
    # )
    # """)
    #
    # db.execute("""
    # CREATE TABLE IF NOT EXISTS reading_sessions (
    #     user TEXT NOT NULL,
    #     book_id INTEGER NOT NULL,
    #     current_page INTEGER NOT NULL,
    #     started_at TIMESTAMP NOT NULL,
    #     last_read_at TIMESTAMP,
    #     PRIMARY KEY (user, book_id),
    #     FOREIGN KEY (user) REFERENCES users(username),
    #     FOREIGN KEY (book_id) REFERENCES books(id)
    # )
    # """)

class UserRepository:
    def __init__(self, db):
        self.db = db

    def save(self, user, password_hash):
        self.db.execute(
            """
            INSERT INTO users (username, display_name, password_hash, created_at)
            VALUES (?, ?, ?, ?)
            """,
            (user.user, user.display_name, password_hash, user.account_created_date)
        )

    def load_all(self):
        rows = self.db.fetchall("SELECT * FROM users")
        return [
            (User(row["username"], row["display_name"]), row["password_hash"])
            for row in rows
        ]

    def load(self, username):
        row = self.db.fetchone(
            "SELECT * FROM users WHERE username = ?",
            (username,)
        )
        if not row:
            return None
        user = User(row["username"], row["display_name"])
        return user, row["password_hash"]

class FriendRepository:
    def __init__(self, db):
        self.db = db

    def add(self, u1, u2):
        self.db.execute(
            "INSERT OR IGNORE INTO friends VALUES (?, ?)",
            (u1, u2)
        )
        self.db.execute(
            "INSERT OR IGNORE INTO friends VALUES (?, ?)",
            (u2, u1)
        )

    def remove(self, u1, u2):
        self.db.execute(
            "DELETE FROM friends WHERE user1 = ? AND user2 = ?",
            (u1, u2)
        )
        self.db.execute(
            "DELETE FROM friends WHERE user1=? AND user2 = ?",
            (u2, u1)
        )

    def load_for_user(self, username):
        rows = self.db.fetchall(
            "SELECT user2 FROM friends WHERE user1 = ?",
            (username, )
        )
        return [row["user2"] for row in rows]

class BookRepository:
    def __init__(self, db):
        self.db = db

    def load_all(self):
        rows = self.db.fetchall("SELECT * FROM books")
        books = []
        for row in rows:
            books.append(Book(
                id = row["id"],
                name = row["name"],
                author = row["author"],
                year = row["year"],
                publisher = row["publisher"],
                language = row["language"],
                total_pages=row["total_pages"]
            ))
        return books

    def get_by_id(self, book_id):
        row = self.db.fetchone(
            "SELECT * FROM books WHERE id = ?",
            (book_id,)
        )
        if not row:
            return None
        return Book(
            id=row["id"],
            name=row["name"],
            author=row["author"],
            year=row["year"],
            publisher=row["publisher"],
            language=row["language"],
            total_pages=row["total_pages"]
        )

class ReadingSessionRepository:
    def __init__(self, db):
        self.db = db

    def save(self, session):
        self.db.execute(
            """
            INSERT OR REPLACE INTO reading_sessions
            (user, book_id, current_page, started_at, last_read_at)
            VALUES (?, ?, ?, ?, ?)
            """,
            (
                session.user.user,
                session.book.id,
                session.current_page,
                session.started_at,
                session.last_read_at
            )
        )

    def delete(self, user, book):
        self.db.execute(
            "DELETE FROM reading_sessions WHERE user = ? AND book_id = ?",
            (user.user, book.id)
        )

class LibraryRepository:
    def __init__(self, db):
        self.db = db

    def create_library(self, name, owner):
        cur = self.db.execute(
            "INSERT INTO libraries (name, owner) VALUES (?, ?)",
            (name, owner)
        )
        return cur.lastrowid

    def add_book(self, library_id, book_id, shelf):
        self.db.execute(
            """
            INSERT OR REPLACE INTO library_books
            (library_id, book_id, shelf)
            VALUES (?, ?, ?)
            """,
            (library_id, book_id, shelf)
        )

    def move_books(self, library_id, book_id, shelf):
        self.db.execute(
            """
            UPDATE library_books
            SET shelf = ?
            WHERE library_id = ? AND book_id = ?
            """,
            (shelf, library_id, book_id)
        )

    def load_libraries_for_user(self, username):
        return self.db.fetchall(
            "SELECT * FROM libraries WHERE owner = ?",
            (username,)
        )

    def load_books_for_library(self, library_id):
        return self.db.fetchall(
            "SELECT * FROM library_books WHERE library_id = ?",
            (library_id,)
        )

class ReviewRepository:
    def __init__(self, db):
        self.db = db

    def save(self, user, book_id, text, rating, created_at):
        self.db.execute(
            """
            INSERT INTO reviews
            (user, book_id, text, rating, created_at)
            VALUES (?, ?, ?, ?, ?)
            """,
            (user, book_id, text, rating, created_at)
        )

    def load_all(self):
        return self.db.fetchall("SELECT * FROM reviews")

class RecommendationsRepository:
    def __init__(self, db):
        self.db = db

    def save(self, from_user, to_user, book_id, message, date):
        self.db.execute(
            """
            INSERT INTO recommendations
            (from_user, to_user, book_id, message, date)
            VALUES (?, ?, ?, ?, ?)
            """,
            (from_user, to_user, book_id, message, date)
        )

    def load_for_user(self, username):
        return self.db.fetchall(
            "SELECT * FROM recommendations WHERE to_user = ?",
            (username,)
        )