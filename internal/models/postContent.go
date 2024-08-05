package models

type PostContent struct {
	Id      string `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content" json:"content"`
}

type PostContentWithFile struct {
	Id               string            `db:"id" json:"id"`
	Title            string            `db:"title" json:"title"`
	Content          string            `db:"content" json:"content"`
	PostContentFiles []PostContentFile `db:"post_content_files" json:"files"`
}
