package models

import "database/sql"

type Post struct {
	Id        string       `db:"id" json:"id"`
	Title     string       `db:"title" json:"title"`
	AuthorId  string       `db:"author_id" json:"author_id"`
	Type      string       `db:"type" json:"type"`
	CreatedAt sql.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}

type PostCreate struct {
	Title string `json:"title"`
}
type PostUpdate struct {
	Title string `json:"title"`
}
