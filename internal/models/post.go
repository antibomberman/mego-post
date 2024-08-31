package models

import (
	"database/sql"
	"time"
)

type Post struct {
	Id        string       `db:"id" json:"id"`
	AuthorId  string       `db:"author_id" json:"author_id"`
	Type      string       `db:"type" json:"type"`
	Image     string       `db:"image" json:"image"`
	CreatedAt sql.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}

type PostDetail struct {
	Id         string               `json:"id"`
	Author     Author               `json:"author"`
	Type       int                  `json:"type"`
	Image      *File                `json:"image"`
	Contents   []PostContentDetails `json:"contents"`
	Categories []CategoryDetails    `json:"categories"`

	CommentCount int `json:"comment_count"`
	LikeCount    int `json:"like_count"`
	ViewCount    int `json:"view_count"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type PostCreate struct {
	AuthorId   string                      `json:"author_id" db:"author_id"`
	Type       int                         `json:"type" db:"type"`
	Image      *FileCreate                 `json:"image" db:"image"`
	Contents   []PostContentCreateOrUpdate `json:"contents" db:"contents"`
	Categories []string                    `json:"categories" db:"categories"`
}
type PostUpdate struct {
	Id         string                      `db:"id" json:"id"`
	Type       int                         `db:"type" json:"type"`
	Image      *FileCreate                 `json:"image" db:"image"`
	Contents   []PostContentCreateOrUpdate `json:"contents"`
	Categories []string                    `json:"categories" db:"categories"`
}
