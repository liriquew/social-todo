-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id        SERIAL PRIMARY KEY,
    username  TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_username ON users (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
