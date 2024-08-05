-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS post_content_files
(
    id              SERIAL PRIMARY KEY,
    post_content_id INT REFERENCES post_contents (id) ,
    filename        TEXT,
    size            BIGINT,
    path             TEXT,
    type            VARCHAR

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS post_content_files;
-- +goose StatementEnd
