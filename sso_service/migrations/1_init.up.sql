CREATE TABLE IF NOT EXISTS users
(
    id        INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_username ON users (username);