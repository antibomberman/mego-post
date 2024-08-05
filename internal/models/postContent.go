package models

type PostContent struct {
	Id      string `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content" json:"content"`
}
type PostContentCreate struct {
	PostId  string `db:"post_id"`
	Title   string `db:"title"`
	Content string `db:"content" json:"content"`
}

type PostContentWithFile struct {
	Id               string            `db:"id" json:"id"`
	Title            string            `db:"title" json:"title"`
	Content          string            `db:"content" json:"content"`
	PostContentFiles []PostContentFile `db:"post_content_files" json:"files"`
}

type PostContentWithFileBinary struct {
	Title            string                  `db:"title" json:"title"`
	Content          string                  `db:"content" json:"content"`
	PostContentFiles []PostContentFileBinary `db:"post_content_files" json:"files"`
}
