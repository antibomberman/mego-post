-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS post_contents
(
    id         SERIAL PRIMARY KEY,
    post_id INT REFERENCES posts (id) ,
    image varchar(255),
    title      TEXT,
    description    TEXT

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS post_contents;
-- +goose StatementEnd
