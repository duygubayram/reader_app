# api.py

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Optional
from passlib.hash import bcrypt
from datetime import datetime
from main import (
    UserManager, User, BookManager,
    ReadingSession, Library, Review
)
from repositories import (
    UserRepository, FriendRepository,
    BookRepository, LibraryRepository,
    ReadingSessionRepository,
    ReviewRepository, RecommendationsRepository,
    Database, init_db
)

app = FastAPI()

# ---------- Setup ----------

db = Database()
init_db(db)

user_repo = UserRepository(db)
friend_repo = FriendRepository(db)
book_repo = BookRepository(db)
reading_repo = ReadingSessionRepository(db)
library_repo = LibraryRepository(db)
review_repo = ReviewRepository(db)
recommendations_repo = RecommendationsRepository(db)

user_manager = UserManager()
book_manager = BookManager()

# ---------- Startup ----------

@app.on_event("startup")
def load_state():
    # users
    for user, _ in user_repo.load_all():
        user_manager.users[user.user] = user

    for username, user in user_manager.users.items():
        friends = friend_repo.load_for_user(username)
        for fname in friends:
            friend = user_manager.users.get(fname)
            if friend:
                user.friends.append(friend)

    # books
    book_manager.load(book_repo.load_all())

    # reading sessions
    rows = db.fetchall("SELECT * FROM reading_sessions")
    for row in rows:
        user = user_manager.users.get(row["user"])
        book = book_manager.get(row["book_id"])
        if user and book:
            session = ReadingSession(user, book)
            session.current_page = row["current_page"]
            user.active_reads[book] = session

    # libraries
    for username, user in user_manager.users.items():
        rows = library_repo.load_libraries_for_user(username)
        for row in rows:
            lib = Library(row["id"], row["name"], user)
            lib_id = row["id"]

            books = library_repo.load_books_for_library(lib_id)
            for b in books:
                book = book_manager.get(b["book_id"])
                lib.shelves[b["shelf"]].append(book)

            user.libraries.append(lib)

    # reviews
    rows = review_repo.load_all()
    for row in rows:
        user = user_manager.users.get(row["user"])
        book = book_manager.get(row["book_id"])
        if user and book:
            review = Review(user, row["text"], row["rating"])
            review.created_at = row["created_at"]
            book.reviews.append(review)

    # recommendations
    for username, user in user_manager.users.items():
        rows = recommendations_repo.load_for_user(username)
        for row in rows:
            from_user = user_manager.users.get(row["from_user"])
            book = book_manager.get(row["book_id"])
            rec = {
                "from": from_user,
                "book": book,
                "message": row["message"],
                "date": row["date"]
            }
            user.recommendations.append(rec)

def hash_password(pw: str) -> str:
    return bcrypt.hash(pw)

def verify_password(pw: str, hashed: str) -> bool:
    return bcrypt.verify(pw, hashed)

def require_user(token: str):
    user = user_manager.users.get(token)
    if not user:
        raise HTTPException(401)
    return user

# ---------- Schemas ----------

# class CreateUserRequest(BaseModel):
#     username: str
#     display_name: str

class RegisterRequest(BaseModel):
    username: str
    display_name: str
    password: str

class LoginRequest(BaseModel):
    username: str
    password: str

class StartReadingRequest(BaseModel):
    username: str
    book_id: int

class TurnPageRequest(BaseModel):
    username: str
    book_id: int
    direction: str
    count: int = 1

class ReviewRequest(BaseModel):
    username: str
    text: str
    rating: int

class RecommendRequest(BaseModel):
    from_user: str
    to_user: str
    book_id: int
    message: Optional[str] = None

# ---------- Users ----------

@app.post("/register")
def register(req: RegisterRequest):
    if user_repo.load(req.username):
        raise HTTPException(400, "User already exists.")

    pw_hash = hash_password(req.password)
    user = User(req.username, req.display_name)
    user_repo.save(user, pw_hash)
    user_manager.users[user.user] = user

    # create default library
    lib_id = library_repo.create_library("My Library", user.user)
    lib = Library(lib_id, "My Library", user)
    user.libraries.append(lib)

    return {
        "status": "created",
        "library_id": lib_id
    }

@app.post("/auth/login")
def login(req: LoginRequest):
    result = user_repo.load(req.username)
    if not result:
        raise HTTPException(401)
    user, pw_hash = result
    if not verify_password(req.password, pw_hash):
        raise HTTPException(401)
    return {"token": user.user} # temp token

# @app.post("/users")
# def create_user(req: CreateUserRequest):
#     try:
#         user = user_manager.create_account(req.username, req.display_name)
#         user_repo.save(user)
#         return {
#             "username": user.user,
#             "display_name": user.display_name
#         }
#     except ValueError as e:
#         raise HTTPException(400, str(e))

@app.get("/me")
def me(token: str):
    user = require_user(token)
    return {
        "username": user.user,
        "display_name": user.display_name,
        "friends": [f.user for f in user.friends],
        "libraries": [lib.name for lib in user.libraries],
    }

@app.get("/users/{username}")
def get_user(username: str):
    user = user_manager.users.get(username)
    if not user:
        raise HTTPException(404)
    return {
        "username": user.user,
        "display_name": user.display_name,
        "friends": [f.user for f in user.friends],
        "libraries": [lib.name for lib in user.libraries]
    }

@app.delete("/users/{username}")
def delete_user(username: str):
    try:
        user_manager.delete_account(username)
        return {"status": "deleted"}
    except ValueError as e:
        raise HTTPException(404, str(e))

# ---------- Friends ----------

@app.post("/users/{u}/friends/{v}")
def add_friend(u: str, v: str):
    u1 = user_manager.users.get(u)
    u2 = user_manager.users.get(v)
    if not u1 or not u2:
        raise HTTPException(404)
    u1.add_friend(u2)
    friend_repo.add(u, v)
    return {"status": "friends"}

@app.delete("/users/{u}/friends/{v}")
def remove_friend(u: str, v: str):
    u1 = user_manager.users.get(u)
    u2 = user_manager.users.get(v)
    if not u1 or not u2:
        raise HTTPException(404)
    u1.remove_friend(u2)
    friend_repo.remove(u, v)
    return {"status": "removed"}

# ---------- Books ----------

@app.get("/books")
def list_books():
    return [
        {
            "id": b.id,
            "name": b.name,
            "author": b.author,
            "year": b.year,
            "language": b.language,
            "publisher": b.publisher,
            "pages": b.total_pages,
            "avg_rating": b.average_rating
        }
        for b in book_manager.all_books()
    ]

@app.get("/books/{book_id}")
def get_book(book_id: int):
    book = book_manager.get(book_id)
    if not book:
        raise HTTPException(404)
    return {
        "id": book.id,
        "name": book.name,
        "author": book.author,
        "year": book.year,
        "language": book.language,
        "publisher": book.publisher,
        "pages": book.total_pages,
        "avg_rating": book.average_rating,
        "reviews": [
            {
                "user": r.reviewer.user,
                "rating": r.rating,
                "text": r.text,
                "likes": r.likes_count
            }
            for r in book.reviews
        ]
    }

# ---------- Reviews ----------

@app.post("/books/{book_id}/reviews")
def review_book(book_id: int, req: ReviewRequest):
    user = user_manager.users.get(req.username)
    book = book_manager.get(book_id)
    if not user or not book:
        raise HTTPException(404)

    review = book.add_review(user, req.text, req.rating)
    # Save the review in the database
    review_repo.save(
        user=user.user,
        book_id=book_id,
        text=req.text,
        rating=req.rating,
        created_at=review.created_at
    )
    return {
        "status": "ok",
        "rating": review.rating
    }

# ---------- Reading ----------

@app.post("/reading/start")
def start_reading(req: StartReadingRequest):
    user = user_manager.users.get(req.username)
    book = book_manager.get(req.book_id)
    if not user or not book:
        raise HTTPException(404)

    session = user.start_reading(book)
    reading_repo.save(session)
    return session.to_database_row()

@app.post("/reading/turn")
def turn_page(req: TurnPageRequest):
    user = user_manager.users.get(req.username)
    book = book_manager.get(req.book_id)
    session = user.active_reads.get(book)
    if not session:
        raise HTTPException(400)

    session.turn_page(req.direction, req.count)
    reading_repo.save(session)
    return session.to_database_row()

@app.post("/reading/stop")
def stop_reading(req: StartReadingRequest):
    user = user_manager.users.get(req.username)
    book = book_manager.get(req.book_id)
    user.stop_reading(book)
    reading_repo.delete(user, book)
    return {"status": "stopped"}

@app.get("/users/{username}/reading")
def list_reading(username: str):
    user = user_manager.users.get(username)
    if not user:
        raise HTTPException(404)
    return [
        s.to_database_row()
        for s in user.active_reads.values()
    ]

# ---------- Libraries ----------

@app.post("/libraries")
def create_library(username: str, name: str):
    user = user_manager.users.get(username)
    if not user:
        raise HTTPException(404)
    lib_id = library_repo.create_library(name, username)
    lib = Library(lib_id, name, user)
    user.libraries.append(lib)
    return {"id": lib_id, "name": name}

@app.post("/libraries/{lib_id}/books")
def add_book(lib_id: int, username: str, book_id: int, shelf: str):
    user = user_manager.users.get(username)
    book = book_manager.get(book_id)
    if not user or not book:
        raise HTTPException(404)
    library = None
    for lib in user.libraries:
        if lib.id == lib_id:
            library = lib
            break
    if not library:
        raise HTTPException(404)
    library_repo.add_book(lib_id, book_id, shelf)
    library.add_book(book, shelf, user)
    return {"status": "added"}

@app.post("/libraries/{lib_id}/move")
def move_book(lib_id: int, username: str, book_id: int, shelf: str):
    user = user_manager.users.get(username)
    book = book_manager.get(book_id)
    if not user or not book:
        raise HTTPException(404)
    library = None
    for lib in user.libraries:
        if lib.id == lib_id:
            library = lib
            break
    if not library:
        raise HTTPException(404)
    library_repo.move_book(lib_id, book_id, shelf)
    for s, books in library.shelves.items():
        if book in books:
            from_shelf = s
            break
    else:
        raise HTTPException(400)
    library.move_book(book, from_shelf, shelf, user)
    return {"status": "moved"}

@app.get("/users/{username}/libraries")
def list_libraries(username: str):
    user = user_manager.users.get(username)
    if not user:
        raise HTTPException(404)
    return [
        {
            "name": lib.name,
            "shelves": {
                k: [b.id for b in v]
                for k, v in lib.shelves.items()
            }
        }
        for lib in user.libraries
    ]

# ---------- Recommendations ----------

@app.post("/recommend")
def recommend(req: RecommendRequest):
    u1 = user_manager.users.get(req.from_user)
    u2 = user_manager.users.get(req.to_user)
    book = book_manager.get(req.book_id)
    if not u1 or not u2 or not book:
        raise HTTPException(404)

    u1.recommend_book(book, u2, req.message)
    recommendations_repo.save(
        req.from_user,
        req.to_user,
        req.book_id,
        req.message,
        datetime.now()
    )
    return {"status": "sent"}

@app.get("/users/{username}/recommendations")
def get_recommendations(username: str):
    user = user_manager.users.get(username)
    if not user:
        raise HTTPException(404)

    return [
        {
            "from": r["from"].user,
            "book": r["book"].name,
            "message": r["message"],
            "date": r["date"]
        }
        for r in user.recommendations
    ]