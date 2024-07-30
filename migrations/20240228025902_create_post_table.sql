-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts
(
    id         SERIAL PRIMARY KEY,
    title      TEXT,
    author_id    INTEGER,
    type    INTEGER,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO posts (user_id, title) VALUES (1, 'Hello World', 'This is my first post');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
-- +goose StatementEnd
