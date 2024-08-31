package models

type PostContent struct {
	Id          string `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Image       string `db:"image" json:"image"`
}

type PostContentDetails struct {
	Id          string `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Image       *File  `db:"image" json:"image"`
}

type PostContentCreateOrUpdate struct {
	PostId      string      `db:"post_id" json:"post_id"`
	Title       string      `db:"title" json:"title"`
	Description string      `db:"description" json:"description"`
	Image       *FileCreate `db:"image" json:"image"`
}
