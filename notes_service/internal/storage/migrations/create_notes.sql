-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    title CHARACTER VARYING(40) UNIQUE NOT NULL,
    note TEXT NOT NULL
    created_at TIMESTAMP NOT NULL,
) 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE 
-- +goose StatementEnd
