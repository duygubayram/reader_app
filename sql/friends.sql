-- bidirectional

CREATE TABLE IF NOT EXISTS friends (
    user1 TEXT NOT NULL,
    user2 TEXT NOT NULL,
    PRIMARY KEY (user1, user2),
    FOREIGN KEY (user1) REFERENCES users(username),
    FOREIGN KEY(user2) REFERENCES users(username)
);