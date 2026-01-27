'''
┌─────────────────────────┐
│  Bubble Tea TUI (Go)    │  ← UI / client
│  - menus                │
│  - keybindings          │
│  - views                │
└───────────▲─────────────┘
            │ HTTP (JSON)
┌───────────┴─────────────┐
│  Python API (FastAPI)   │  ← backend
│  - users                │
│  - libraries            │
│  - reading sessions     │
│  - SQL persistence      │
└───────────▲─────────────┘
            │
┌───────────┴─────────────┐
│  SQLite / Postgres      │
└─────────────────────────┘

==============================================================================

Domain Model Overview

User
 ├── user: str                      # unique username (primary identifier)
 ├── display_name: str
 ├── account_created_date: datetime
 ├── friends: list[User]            # bidirectional friendships (in-memory)
 ├── libraries: list[Library]       # user-owned libraries
 ├── recommendations: list[dict]    # incoming book recommendations
 ├── active_reads: dict[Book, ReadingSession]
 │     └── tracks active reading sessions per book

UserManager
 ├── users: dict[str, User]         # username → User
 ├── creates and deletes accounts
 ├── prepares User data for persistence (SQL-ready)

Library
 ├── name: str
 ├── owner: User
 ├── shelves: dict[str, list[Book]]
 │     ├── "to_read"
 │     ├── "currently_reading"
 │     └── "read"
 ├── enforces owner-only modifications
 ├── SQL-mappable via to_database_row()

ReadingSession
 ├── user: User
 ├── book: Book
 ├── current_page: int
 ├── started_at: datetime
 ├── last_read_at: datetime | None
 ├── auto-finishes when last page is reached
 ├── SQL-mappable via to_database_row()

Book
 ├── id: int                        # unique identifier (database-backed)
 ├── name: str
 ├── author: str
 ├── language: str
 ├── year: int
 ├── publisher: str
 ├── total_pages: int
 ├── reviews: list[Review]
 ├── equality & hashing based on id

BookManager
 ├── books: dict[int, Book]         # book_id → Book
 ├── loads books from external database
 ├── acts as single source of truth for Book instances

Review
 ├── reviewer: User
 ├── rating: int (1–5)
 ├── text: str
 ├── created_at: datetime
 ├── likes: set[User]               # prevents duplicate likes
'''

from datetime import datetime

class UserManager:
    def __init__(self):
        self.users = {} # username -> user

    def create_account(self, username, display_name):
        if username in self.users:
            raise ValueError("Username already exists.")
        user = User(username, display_name)
        self.users[username] = user
        return user

    def delete_account(self, username):
        user = self.users.pop(username, None)
        if not user:
            raise ValueError("User not found.")
        for friend in list(user.friends): # iterate over copy not the actual list
            friend.friends.remove(user)

    def to_database_row(self, user):
        return {
            "username": user.user,
            "display_name": user.display_name,
            "created_at": user.account_created_date
        }

class User:
    def __init__(self, user, display_name):
        self.user = user # refers to unique username
        self.display_name = display_name
        self.account_created_date = datetime.now()
        self.friends = []
        self.libraries = []
        self.recommendations = []
        self.active_reads = {}
        self.create_library("My Library")

    def add_friend(self, friend_request):
        if friend_request not in self.friends and friend_request != self:
            self.friends.append(friend_request)
            friend_request.friends.append(self)

    def remove_friend(self, friend):
        if friend in self.friends:
            self.friends.remove(friend)
            friend.friends.remove(self)

    def recommend_book(self, book, friend, message = None):
        if friend not in self.friends:
            raise ValueError("Can only recommend books to friends.")
        recommendation = {
            "from": self,
            "book": book,
            "message": message,
            "date": datetime.now()
        }
        friend.receive_recommendation(recommendation)

    def receive_recommendation(self, recommendation):
        self.recommendations.append(recommendation)

    def create_library(self, name):
        library = Library(name, self)
        self.libraries.append(library)
        return library

    @property
    def primary_library(self):
        return self.libraries[0]

    def start_reading(self, book):
        if book in self.active_reads:
            return self.active_reads[book]
        library = self.primary_library
        for shelf in library.shelves:
            if book in library.shelves[shelf]:
                library.move_book(book, shelf, "currently_reading", self)
                break
        session = ReadingSession(self, book)
        self.active_reads[book] = session
        return session

    def stop_reading(self, book):
        session = self.active_reads.pop(book, None)
        if not session:
            return
        if session.current_page >= book.total_pages:
            library = self.primary_library
            library.move_book(book, "currently_reading", "read", self)

    def __repr__(self):
        return f"User({self.user})"

    #@classmethod
    #def add_user(cls):
    #    cls.num_of_users += 1

class ReadingSession:
    def __init__(self, user, book):
        self.user = user
        self.book = book
        self.current_page = 1
        self.started_at = datetime.now()
        self.last_read_at = None

    def turn_page(self, direction, count=1):
        if direction == "forward":
            self.current_page = min(self.book.total_pages, self.current_page + count)
        elif direction == "back":
            self.current_page = max(1, self.current_page - count)
        if self.current_page >= self.book.total_pages:
            self.user.stop_reading(self.book)
        self.last_read_at = datetime.now()

    #def is_finished(self):
    #    return self.current_page >= self.book.total_pages

    def to_database_row(self):
        return{
            "user": self.user.user,
            "book_id": self.book.id,
            "current_page": self.current_page,
            "started_at": self.started_at,
            "last_read_at": self.last_read_at
        }

# incase they want multiple libraries (eg personal, school, work)
class Library:
    permission_error = "Only the owner can modify this library."

    def __init__(self, name, owner):
        self.name = name
        self.owner = owner
        self.shelves = {
            "to_read": [],
            "currently_reading": [],
            "read": []
        }

    def add_book(self, book, shelf, user):
        if user != self.owner:
            raise PermissionError(self.permission_error)
        if shelf in self.shelves and book not in self.shelves[shelf]:
            self.shelves[shelf].append(book)

    def remove_book(self, book, user):
        if user != self.owner:
            raise PermissionError(self.permission_error)
        for shelf in self.shelves.values():
            while book in shelf:
                shelf.remove(book)

    def move_book(self, book, from_shelf, to_shelf, user):
        if user != self.owner:
            raise PermissionError(self.permission_error)
        if from_shelf not in self.shelves or to_shelf not in self.shelves:
            raise ValueError("Invalid shelf name.")
        if book in self.shelves[from_shelf]:
            self.shelves[from_shelf].remove(book)
        if book not in self.shelves[to_shelf]:
            self.shelves[to_shelf].append(book)

    def create_shelf(self, shelf_name, user):
        if user != self.owner:
            raise PermissionError(self.permission_error)
        if shelf_name in self.shelves:
            raise ValueError("Shelf already exists.")
        self.shelves[shelf_name] = []

    def to_database_row(self):
        return {
            "name": self.name,
            "owner": self.owner.user
        }

class BookManager:
    def __init__(self):
        self.books = {} # by id

    def load(self, books):
        for book in books:
            self.books[book.id] = book

    def load_from_database(self, rows):
        # rows = iterable of dicts or tuples from sql
        for row in rows:
            book = Book(
                id = row["id"],
                name = row["name"],
                author = row["author"],
                language = row["language"],
                year = row["year"],
                publisher = row["publisher"],
                total_pages=row["total_pages"]
            )
            self.books[book.id] = book

    def get(self, book_id):
        return self.books.get(book_id)

    def all_books(self):
        return list(self.books.values())

class Book:
    def __init__(self, id, name, author, year, publisher, language, total_pages):
        self.id = id
        self.name = name
        self.author = author
        self.language = language
        self.year = year
        self.publisher = publisher
        self.total_pages = total_pages
        self.reviews = []

    def __eq__(self, other):
        return isinstance(other, Book) and self.id == other.id

    def __hash__(self):
        return hash(self.id)

    def add_review(self, user, text, rating):
        for review in self.reviews:
            if review.reviewer == user:
                raise ValueError("User has already reviewed this book.")
        review = Review(user, text, rating)
        self.reviews.append(review)
        return review

    @property
    def average_rating(self):
        if not self.reviews:
            return None
        return sum(review.rating for review in self.reviews) / len(self.reviews)

    #def add_rating(self, stars_count):
    #    if 1 <= stars <= 5:
    #        self.ratings.append(stars)
    def __repr__(self):
        return f"Book(id={self.id}, name='{self.name}')"

class Review:
    def __init__(self, reviewer, text, rating):
        if not 1 <= rating <= 5:
            raise ValueError("Rating must be between 1 and 5.")
        self.reviewer = reviewer
        self.text = text
        self.rating = rating
        self.created_at = datetime.now()
        self.likes = set() # prevents duplicate likes from same user

    def like(self, user):
        if user == self.reviewer:
            raise ValueError("Cannot like your own review")
        self.likes.add(user)

    def unlike(self, user):
        self.likes.discard(user)

    @property
    def likes_count(self):
        return len(self.likes)