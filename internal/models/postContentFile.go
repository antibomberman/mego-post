package models

type PostContentFile struct {
	Id            string `db:"id" json:"id"`
	PostContentId string `db:"post_content_id" json:"-"`
	FileName      string `db:"file_name" json:"filename"`
	ContentType   string `db:"content_type" json:"content_type"`
	Url           string `db:"url" json:"url"`
}
type PostContentFileCreate struct {
	PostContentId string `db:"post_content_id" json:"-"`
	FileName      string `db:"file_name" json:"filename"`
	ContentType   string `db:"content_type" json:"content_type"`
	Data          []byte `db:"data" json:"-"`
}
