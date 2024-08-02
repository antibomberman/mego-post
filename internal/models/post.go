package models

type Post struct {
	Id        string `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	AuthorId  string `db:"author_id" json:"author_id"`
	Type      string `db:"type" json:"type"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	DeletedAt string `db:"deleted_at" json:"deleted_at"`
}
type PostDetails struct {
	Id        string        `json:"id"`
	Title     string        `json:"title"`
	Author    Author        `json:"author"`
	Type      string        `json:"type"`
	Contents  []PostContent `json:"contents"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	DeletedAt string        `json:"deleted_at"`
}

type PostCreate struct {
	Title string `json:"title"`
}
type PostUpdate struct {
	Title string `json:"title"`
}
