-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS hide_posts
(
    id      SERIAL PRIMARY KEY,
    user_id INT,
    post_id INT REFERENCES post_contents (id)

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS hide_posts;
-- +goose StatementEnd
