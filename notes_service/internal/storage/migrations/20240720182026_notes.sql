-- +goose Up
-- +goose StatementBegin

/*
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    title CHARACTER VARYING(40) UNIQUE NOT NULL,
    note TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
) 
*/
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX unique_title_per_owner ON notes(owner_id, title);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notes
-- +goose StatementEnd
