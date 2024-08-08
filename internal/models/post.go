package models

import (
	"database/sql"
	"time"
)

type Post struct {
	Id        string       `db:"id" json:"id"`
	Title     string       `db:"title" json:"title"`
	AuthorId  string       `db:"author_id" json:"author_id"`
	Type      string       `db:"type" json:"type"`
	CreatedAt sql.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}

type PostDetail struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Author Author `json:"author"`
	Type   string `json:"type"`

	Contents []PostContentWithFile `json:"contents"`

	CommentCount int `json:"comment_count"`
	LikeCount    int `json:"like_count"`
	RepostCount  int `json:"repost_count"`
	ViewCount    int `json:"view_count"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type PostCreate struct {
	Title    string                      `json:"title"`
	AuthorId string                      `json:"author_id"`
	Type     string                      `json:"type"`
	Contents []PostContentCreateOrUpdate `json:"contents"`
}
type PostUpdate struct {
	Id       string                      `db:"id" json:"id"`
	Title    string                      `db:"title" json:"title"`
	Type     string                      `db:"type" json:"type"`
	Contents []PostContentCreateOrUpdate `json:"contents"`
}
