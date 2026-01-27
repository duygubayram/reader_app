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
    with open("books.sql", "r", encoding="utf-8") as f:
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

    def save(self, user):
        self.db.execute(
            """
            INSERT INTO users (username, display_name, created_at)
            VALUES (?, ?, ?)
            """,
            (user.user, user.display_name, user.account_created_date)
        )

    def load_all(self):
        rows = self.db.fetchall("SELECT * FROM users")
        return [User(row["username"], row["display_name"]) for row in rows]

    def load(self, username):
        row = self.db.fetchone(
            "SELECT * FROM users WHERE username = ?",
            (username,)
        )
        if not row:
            return None
        return User(row["username"], row["display_name"])

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