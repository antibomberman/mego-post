package models

type PostContent struct {
	Id          string `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	File        File   `db:"file" json:"file"`
}
type PostContentCreate struct {
	PostId      string     `db:"post_id" json:"post_id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	File        FileCreate `db:"file" json:"file"`
}

type PostContentCreateOrUpdate struct {
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	File        FileCreate `db:"file" json:"file"`
}
