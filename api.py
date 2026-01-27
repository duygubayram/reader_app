# api.py

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Optional
from main import (
    UserManager, User, BookManager,
    ReadingSession, Library
)
from repositories import (
    UserRepository, BookRepository,
    ReadingSessionRepository,
    Database, init_db
)

app = FastAPI()

# ---------- Setup ----------

db = Database()
init_db(db)

user_repo = UserRepository(db)
book_repo = BookRepository(db)
reading_repo = ReadingSessionRepository(db)

user_manager = UserManager()
book_manager = BookManager()

# ---------- Startup ----------

@app.on_event("startup")
def load_state():
    # users
    for user in user_repo.load_all():
        user_manager.users[user.user] = user

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

# ---------- Schemas ----------

class CreateUserRequest(BaseModel):
    username: str
    display_name: str

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

@app.post("/users")
def create_user(req: CreateUserRequest):
    try:
        user = user_manager.create_account(req.username, req.display_name)
        user_repo.save(user)
        return {
            "username": user.user,
            "display_name": user.display_name
        }
    except ValueError as e:
        raise HTTPException(400, str(e))

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
    return {"status": "friends"}

@app.delete("/users/{u}/friends/{v}")
def remove_friend(u: str, v: str):
    u1 = user_manager.users.get(u)
    u2 = user_manager.users.get(v)
    if not u1 or not u2:
        raise HTTPException(404)
    u1.remove_friend(u2)
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