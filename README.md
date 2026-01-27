# reader_app
oop practice with sql, api and tui -- mostly to learn app building after backend processes

## Architecture Overview

```
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
```

---

## Domain Model Overview

```
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
```
