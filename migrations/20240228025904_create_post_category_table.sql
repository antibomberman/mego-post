-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS post_categories
(
    id         SERIAL PRIMARY KEY,
    post_id INT REFERENCES posts (id) ,
    category_id INT REFERENCES categories (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS post_categories;
-- +goose StatementEnd
